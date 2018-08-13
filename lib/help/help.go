package help

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/meeDamian/bc1toolkit/lib/common"
	"github.com/sirupsen/logrus"
)

type Opts struct {
	Version  bool     `long:"version" short:"v" description:"Show version and exit"`
	Verbose  []bool   `long:"verbose" short:"V" description:"Enable verbose logging. Specify twice to increase verbosity"`
	TorMode  string   `long:"tor-mode" description:"When to use Tor. \"native\" - .onion addresses only. \"auto\" - see above for details." choice:"always" choice:"auto" choice:"native" choice:"never" default:"auto"`
	TorSocks []string `long:"tor" description:"\"host:port\" to Tor's SOCKS proxy" default:"localhost:9050" default:"localhost:9150" default-mask:"localhost:9050 or localhost:9150"`
}

var (
	version,
	buildStamp,
	gitHash,
	builder string

	opts Opts

	parser = flags.NewParser(&opts, flags.Default)
)

func Customize(usage, description, torAutoBehaviour, name string, data interface{}) {
	if usage != "" {
		parser.Usage = usage
	}

	if description != "" {
		if torAutoBehaviour != "" {
			description = fmt.Sprintf("%s\n\nTor \"auto\" behaviour: %s", description, torAutoBehaviour)
		}

		parser.LongDescription = description
	}

	if name != "" {
		parser.AddGroup(name, "", data)
	}
}

func Parse() ([]string, Opts) {
	args, err := parser.Parse()
	if err != nil {
		// show help message and exit
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}

		fmt.Println("Unable to parse flags:", err.Error())
		os.Exit(1)
	}

	if opts.Version {
		fmt.Printf("%-16s: %s\n%-16s: %s\n%-16s: %s\n%-16s: %s\n",
			"Version", version,
			"Git Commit Hash", gitHash,
			"UTC built time", buildStamp,
			"Built by", builder,
		)
		os.Exit(0)
	}

	switch len(opts.Verbose) {
	case 0:
		common.Logger.SetLevel(logrus.WarnLevel)

	case 1:
		common.Logger.SetLevel(logrus.InfoLevel)

	case 2:
		fallthrough
	default: // more than 2
		common.Logger.SetLevel(logrus.DebugLevel)
	}

	return args, opts
}
