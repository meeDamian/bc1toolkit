// +build !windows,!darwin

package common

import (
	"os"
	"path/filepath"
)

var cacheBase string

func init() {
	if os.Getenv("XDG_CACHE_HOME") != "" {
		cacheBase = os.Getenv("XDG_CACHE_HOME")
		return
	}

	cacheBase = filepath.Join(os.Getenv("HOME"), ".cache")
}
