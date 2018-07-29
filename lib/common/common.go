package common

import (
	"github.com/meeDamian/bc1toolkit/lib/connstring"
	"github.com/meeDamian/bc1toolkit/lib/tor"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
)

var Log = logrus.New()

type Dialers struct {
	Tor      proxy.Dialer
	ClearNet proxy.Dialer
	mode     string
}

type SpeakFn func(dialer proxy.Dialer, addr connstring.ConnString, testNet bool) (interface{}, error)

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
