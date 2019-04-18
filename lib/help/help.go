package help

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/meeDamian/bc1toolkit/lib/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const DisableTor = ""

type Opts struct {
	Version bool   `long:"version" short:"v" no-ini:"yes" description:"Show version and exit"`
	Verbose []bool `long:"verbose" short:"V" description:"Enable verbose logging. Specify twice to increase verbosity"`

	// INI config
	Config string `long:"config" no-ini:"yes" description:"Use config from file.  CLI flags take precedence." default-mask:"./bc1toolkit.conf"`
	Save   bool   `long:"save" no-ini:"yes" description:"Run and update config file with current options"`

	// Tor
	TorMode  string   `long:"tor-mode" description:"When to use Tor. \"native\" - end-to-end .onion only. \"auto\" - see above for details." choice:"always" choice:"auto" choice:"native" choice:"never" default:"auto"`
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

	if torAutoBehaviour == DisableTor {
		parser.FindOptionByLongName("tor").Hidden = true
		parser.FindOptionByLongName("tor-mode").Hidden = true
		parser.FindOptionByLongName("tor-mode").Choices = []string{}
	} else if description != "" {
		description = fmt.Sprintf("%s\n\nTor \"auto\" behaviour: %s", description, torAutoBehaviour)
	}

	if description != "" {
		parser.LongDescription = description
	}

	if name != "" {
		_, err := parser.AddGroup(name, "", data)
		if err != nil {
			fmt.Printf("couldn't add %s group %v", name, err)
			os.Exit(1)
		}
	}
}

func processConfigFile(fileName string, bestEffortRead bool) error {
	iniParser := flags.NewIniParser(parser)
	err := iniParser.ParseFile(fileName)
	if err != nil {
		if bestEffortRead {
			return nil
		}

		return errors.Wrapf(err, "Unable to read config")
	}

	common.Logger.Get().WithField("config-file", fileName).Debug("Config read successfully")

	// Parse CLI again, because CLI takes precedent above config from file
	_, _ = parser.Parse()

	if opts.Save && !bestEffortRead {
		err = iniParser.WriteFile(opts.Config, flags.IniNone)
		if err != nil {
			return errors.Wrap(err, "Unable to save config")
		}
	}

	return nil
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

	if len(opts.Config) == 0 && opts.Save {
		fmt.Printf("To --save config, you need to pass path to it explicitly, ex:\n\t--config=./bc1toolkit.conf --save\n\n")
		os.Exit(1)
	}

	err = processConfigFile(opts.Config, len(opts.Config) == 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
