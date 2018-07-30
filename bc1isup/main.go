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

const description = `Checks addresses for running Bitcoin nodes. When addresses are both piped-in and provided at command line, piped ones are first.

Each address provided outputs its own line with corresponding node status info (customizable with --output=?).
Exit code of 0 is returned only if each address provided had at least one running node.`

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
	Opts struct {
		TestNet bool   `long:"testnet" short:"T" description:"Check for testnet node"`
		MainNet bool   `long:"mainnet" short:"M" description:"Check for mainnet node"`
		AutoNet bool   `no-flag:"can be used to determine if check was requested or is an auto-fallback"`
		Output  string `long:"output" short:"o" description:"Choose line format: 'json' for JSON array. 'simple' for a single \"up\" or \"down\". 'none' for no output, and only exit code" default:"json" choice:"json" choice:"simple" choice:"none"`
	}

	connStrings []string

	opts help.Opts
)

func init() {
	help.Customize(
		"[OPTIONS] (domain|IP)[:port] ...",
		description,
		"bc1isup", &Opts,
	)

	// read parameters passed to a binary
	connStrings, opts = help.Parse()

	// check for stuff being piped-in
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		rawStdin, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// process input, split into addresses and strip possible quotes
		var stdinConnStrings []string
		for _, a := range strings.Split(strings.TrimSpace(string(rawStdin)), "\n") {
			stdinConnStrings = append(stdinConnStrings, strings.Replace(a, "\"", "", -1))
		}

		// pipe data first, and then command ones seems more natural
		connStrings = append(stdinConnStrings, connStrings...)
	}

	if len(connStrings) < 1 {
		fmt.Println("At least one node conn-string needs to be provided")
		os.Exit(1)
	}
}

func attemptCommunication(explicitlyRequested, testNet bool, dialer proxy.Dialer, c connstring.ConnString) (version interface{}) {
	if !explicitlyRequested && !Opts.AutoNet {
		return nil
	}

	version, err := btc.Speak(dialer, c, testNet)
	if err != nil {
		if Opts.AutoNet {
			common.Log.Debugln(err)
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

	version := attemptCommunication(Opts.MainNet, false, dialer, c)
	if version != nil {
		found = append(found, version)
	}

	version = attemptCommunication(Opts.TestNet, true, dialer, c)
	if version != nil {
		found = append(found, version)
	}

	return
}

func main() {
	onlyLocal := true
	noTor := true

	var cs []connstring.ConnString
	for _, c := range connStrings {
		conn, err := connstring.Parse(c)
		if err != nil {
			fmt.Printf("%s is not valid: %v\n", c, err)
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
	if !Opts.TestNet && !Opts.MainNet {
		Opts.AutoNet = true
	}

	// skip Tor altogether when possible
	if onlyLocal {
		opts.TorMode = "never"

	} else if opts.TorMode == "native" && noTor {
		opts.TorMode = "never"
	}

	dialers, err := common.GetDialers(opts.TorMode, opts.TorSocks)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(len(cs) + 1) // all checks + the final results aggregation goroutine

	results := make(chan result, len(cs))

	// launch all requested checks in parallel
	for id, c := range cs {
		go func(id int, c connstring.ConnString) {
			found, err := checkConnString(dialers, c)
			if err != nil {
				found = []interface{}{nodeError{
					c.Host,
					err.Error(),
				}}
			}

			results <- result{id, found}

			wg.Done()
		}(id, c)
	}

	// get all checks and output them in the same order as provided
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

					if len(item) == 0 {
						// TODO: add nothing-found error(?)
					}

					v, err := json.Marshal(item)
					if err != nil {
						common.Log.Errorln("unable to marshall response…")
					}

					fmt.Println(string(v))
					delete(received, last)
				}
			}

			if last >= max {
				close(results)
			}
		}

		wg.Done()
	}()

	wg.Wait()
}
