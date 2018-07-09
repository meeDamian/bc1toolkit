package common

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

var (
	Opts struct {
		Version       bool `short:"v" long:"version" description:"Show version of the binary"`
		HumanReadable bool `short:"H" long:"human-readable" description:"Show in a human-readable format"`
	}

	buildStamp,
	gitHash,
	builder string

	Parser *flags.Parser
)

func init() {
	Parser = flags.NewParser(&Opts, flags.Default)
}

func DefaultActions() {
	if Opts.Version {
		fmt.Printf("Git Commit Hash: %s\nUTC build time : %s\nBuilt by       : %s\n", gitHash, buildStamp, builder)
		os.Exit(0)
	}
}
