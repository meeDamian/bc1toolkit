package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/meeDamian/bc1toolkit/lib/btc"
	"github.com/meeDamian/bc1toolkit/lib/common"
	"github.com/meeDamian/bc1toolkit/lib/connstring"
	"github.com/meeDamian/bc1toolkit/lib/help"
	"golang.org/x/net/proxy"
)

const (
	BinaryName = "bc1isup"

	description = `Checks addresses for running Bitcoin nodes. When addresses are both piped-in and provided at command line, piped ones are first.

Each address provided, outputs its own line with corresponding node status info (customizable with --output=?).
Exit code of 0 is returned only if each address provided had at least one running node.`

	torBehaviour = `tries using Tor, if not available, falls back to clearnet.`
)

type (
	nodeError struct {
		Address string `json:"address"`
		Error   string `json:"error"`
	}

	result struct {
		Id  int
		Out []interface{}
	}
)

var (
	opts struct {
		TestNet bool   `long:"testnet" short:"T" description:"Check for testnet node"`
		MainNet bool   `long:"mainnet" short:"M" description:"Check for mainnet node"`
		AutoNet bool   `no-flag:"can be used to determine if check was requested or is an auto-fallback"`
		Output  string `long:"output" short:"o" description:"Choose line format: 'json' for JSON array. 'simple' for a single \"up\" or \"down\". 'none' for no output, and only exit code" default:"json" choice:"json" choice:"simple" choice:"none"`
	}
	commonOpts help.Opts

	addresses []string
)

// NOTE: all errors returned here are quoted strings to preserve `jq` compatibility
func init() {
	common.Logger.Name(BinaryName)

	help.Customize(
		"[OPTIONS] (domain|IP)[:port] ...",
		description,
		torBehaviour,
		BinaryName, &opts,
	)

	// read parameters passed to a binary
	addresses, commonOpts = help.Parse()

	// check for stuff being piped-in
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		rawStdin, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Printf(`"%s"\n`, err)
			os.Exit(1)
		}

		// process input, split into addresses and trim possible whitespaces & quotes
		var stdinAddresses []string
		for _, a := range strings.Split(string(rawStdin), "\n") {
			// if there are extra whitespaces in the source file
			trimmedAddress := strings.TrimSpace(a)

			// if addresses are piped from `jq` w/o `-r` provided
			trimmedAddress = strings.Trim(trimmedAddress, "\"")

			stdinAddresses = append(stdinAddresses, trimmedAddress)
		}

		// pipe data first, and then command ones seems more natural
		addresses = append(stdinAddresses, addresses...)
	}

	if len(addresses) < 1 {
		fmt.Println(`"At least one IP address or hostname needs to be provided"`)
		os.Exit(1)
	}
}

func attemptCommunication(explicitlyRequested, testNet bool, dialer proxy.Dialer, c connstring.ConnString) (version interface{}) {
	if !explicitlyRequested && !opts.AutoNet {
		return nil
	}

	version, err := btc.Speak(dialer, c, testNet)
	if err != nil {
		if opts.AutoNet {
			common.Logger.Get().Debugln(err)
			return nil
		}

		return nodeError{c.Raw, err.Error()}
	}

	return version

}

func checkConnString(dialers common.Dialers, c connstring.ConnString) (found []interface{}, err error) {
	dialer, err := dialers.Default(c.IsTor(), c.Local)
	if err != nil {
		return nil, err
	}

	version := attemptCommunication(opts.MainNet, false, dialer, c)
	if version != nil {
		found = append(found, version)
	}

	version = attemptCommunication(opts.TestNet, true, dialer, c)
	if version != nil {
		found = append(found, version)
	}

	return
}

func main() {
	onlyLocal, noTor := true, true

	var cs []connstring.ConnString
	for _, c := range addresses {
		conn, err := connstring.Parse(c)
		if err != nil {
			fmt.Printf(`"%s is not valid: %v"\n`, c, err)
			os.Exit(1)
		}

		if !conn.Local {
			onlyLocal = false
		}

		if conn.IsTor() {
			noTor = false
		}

		cs = append(cs, conn)
	}

	// if neither is specified, perform auto check
	if !opts.TestNet && !opts.MainNet {
		opts.AutoNet = true
	}

	// skip Tor altogether when possible
	if onlyLocal {
		common.Logger.Get().Debugln("only local addresses provided: disabling To completely")
		commonOpts.TorMode = "never"

	} else if commonOpts.TorMode == "native" && noTor {
		common.Logger.Get().Debugln("--tor-mode=native set and no Tor addresses provided: disabling To completely")
		commonOpts.TorMode = "never"
	}

	// Return only dialers that will be used in requests
	dialers, err := common.GetDialers(commonOpts.TorMode, commonOpts.TorSocks)
	if err != nil {
		fmt.Printf(`"%s"\n`, err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(len(cs) + 1) // all checks + the final results aggregation goroutine

	results := make(chan result, len(cs))

	exitCode := 0

	// launch all requested checks in parallel
	for id, c := range cs {
		go func(id int, c connstring.ConnString) {
			found, err := checkConnString(dialers, c)
			if err != nil {
				found = []interface{}{nodeError{
					c.Host,
					err.Error(),
				}}

			} else if found == nil {
				exitCode = 1

				found = []interface{}{}
			}

			for _, x := range found {
				if _, ok := x.(nodeError); ok {
					exitCode = 1
				}
			}

			results <- result{id, found}

			wg.Done()
		}(id, c)
	}

	// receive all checks and output them in the same order as provided
	go func() {
		received := make(map[int][]interface{})
		last, max := 0, len(cs)

		for r := range results {
			received[r.Id] = r.Out

			if r.Id >= last {
				for ; last <= max; last++ {
					item, ok := received[last]
					if !ok {
						break
					}

					delete(received, last)

					// output is irrelevant. Just wait until all are received
					if opts.Output == "none" {
						continue
					}

					if opts.Output == "simple" {
						if len(item) == 0 {
							fmt.Println("down")
							continue
						}

						out := "up"
						for _, x := range item {
							if _, ok := x.(nodeError); !ok {
								continue
							}
							out = "down"
						}
						fmt.Println(out)
					}

					if opts.Output == "json" {
						v, err := json.Marshal(item)
						if err != nil {
							exitCode = 1

							common.Logger.Get().Errorf("unable to marshall response: %#v", item)
							fmt.Println(`[{"error": "unable to marshall response"}]`)
							continue
						}

						fmt.Println(string(v))
					}
				}
			}

			if last >= max {
				close(results)
			}
		}

		wg.Done()
	}()

	wg.Wait()
	os.Exit(exitCode)
}
