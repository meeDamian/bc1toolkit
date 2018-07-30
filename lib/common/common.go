package common

import (
	"path"

	"github.com/meeDamian/bc1toolkit/lib/tor"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
)

const cacheDir = "com.meedamian.bc1toolkit"

type Dialers struct {
	Tor      proxy.Dialer
	ClearNet proxy.Dialer
	mode     string
}

// TODO: deprecated
var Log = logrus.New()

func (d Dialers) Default(isTor, isLocal bool) (proxy.Dialer, error) {
	if isLocal {
		return d.ClearNet, nil
	}

	if d.mode == "never" && isTor {
		return nil, errors.New("can't check .onion addresses with --tor-mode=never")
	}

	if d.mode == "native" && isTor {
		if d.Tor == nil {
			return nil, errors.New("can't check .onion addresses with Tor unavailable")
		}

		return d.Tor, nil
	}

	if d.mode == "auto" && d.Tor != nil {
		return d.Tor, nil
	}

	if d.mode == "always" {
		return d.Tor, nil
	}

	return d.ClearNet, nil
}

func GetLogger(binaryName string) *logrus.Entry {
	return logrus.New().WithField("binary", binaryName)
}

func GetDialers(torMode string, torSocks []string) (d Dialers, _ error) {
	d = Dialers{
		mode:     torMode,
		ClearNet: proxy.Direct,
	}

	if torMode == "never" {
		return
	}

	torDialer, err := tor.GetWorkingTor(torSocks)
	if err != nil {
		if torMode == "always" {
			return Dialers{}, errors.New("can't connect to Tor & --tor-mode=always set")
		}

		Log.Debugln("Tor connection not available, using clearnet only.", err)
	}

	d.Tor = torDialer
	return
}

func GetCacheDir() string {
	return path.Join(cacheBase, cacheDir)
}
