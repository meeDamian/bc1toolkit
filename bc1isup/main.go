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
	"github.com/meeDamian/bc1toolkit/lib/ln"
	"golang.org/x/net/proxy"
)

// TODO: LN nodes
// TODO: minimise requests
// TODO: README
// TODO: cache LN pubkeys per node

const description = `Check conn-strings for running Bitcoin and LN nodes.`

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
		Network     string `long:"network" short:"N" description:"Check against specified network only" choice:"btc" choice:"ln" choice:"all" default:"all"`
		TestNet     bool   `long:"testnet" short:"T" description:"Check for testnet"`
		MainNet     bool   `long:"mainnet" short:"M" description:"Check for mainnet"`
		AutoNet     bool   `no-flag:"can be used to determine if check was requested or is an auto-fallback"`
		PubKeyCache bool   `long:"no-pubkey-cache" description:"Do not cache Lightning Network pubkey for future runs. Default cache location is ~/.cache/bc1toolkit/ln-nodes"`
		Output      string `long:"output" short:"o" description:"Output mode. 'json' returns JSON array for each input address. 'simple' returns up/down per input address. 'none' returns no output, but returns exit code 0 if all connstrings were up and 1 if any failed." default:"json" choice:"json" choice:"simple" choice:"none"`
	}

	connStrings []string

	opts help.Opts
)

func init() {
	help.Customize(
		"[OPTIONS] [pubkey@](domain|IP)[:port] ...",
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

func tryToConverse(fn common.SpeakFn, dialer proxy.Dialer, c connstring.ConnString, testNet bool) (interface{}, error) {
	version, err := fn(dialer, c, testNet)
	if err != nil {
		if Opts.AutoNet {
			common.Log.Debugln(err)
			return nil, nil
		}

		return nil, err
	}

	return version, nil
}

func infoOrError(explicitlyRequested, testNet bool, fn common.SpeakFn, dialer proxy.Dialer, c connstring.ConnString) (version interface{}) {
	if !explicitlyRequested && !Opts.AutoNet {
		return nil
	}

	version, err := tryToConverse(fn, dialer, c, testNet)
	if err != nil {
		return nodeError{c.Raw, err.Error()}
	}

	return version

}

func checkConnString(dialers common.Dialers, c connstring.ConnString) (found []interface{}, err error) {
	dialer, err := dialers.Default(c.IsTor(), c.Local)
	if err != nil {
		return nil, err
	}

	if Opts.Network == "btc" || Opts.Network == "all" {
		version := infoOrError(Opts.MainNet, false, btc.Speak, dialer, c)
		if version != nil {
			found = append(found, version)
		}

		version = infoOrError(Opts.TestNet, true, btc.Speak, dialer, c)
		if version != nil {
			found = append(found, version)
		}
	}

	if Opts.Network == "ln" || Opts.Network == "all" {
		version := infoOrError(Opts.MainNet, false, ln.Speak, dialer, c)
		if version != nil {
			found = append(found, version)
		}

		version = infoOrError(Opts.TestNet, true, ln.Speak, dialer, c)
		if version != nil {
			found = append(found, version)
		}
	}

	return
}

func main() {
	var cs []connstring.ConnString
	for _, c := range connStrings {
		conn, err := connstring.Parse(c)
		if err != nil {
			fmt.Printf("%s is not valid: %v\n", c, err)
			os.Exit(1)
		}

		cs = append(cs, conn)
	}

	// if neither is specified, perform auto check
	if !Opts.TestNet && !Opts.MainNet {
		Opts.AutoNet = true
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
