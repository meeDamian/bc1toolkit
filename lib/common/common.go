package common

import (
	"path"

	"github.com/meeDamian/bc1toolkit/lib/tor"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
)

const (
	cacheDir      = "com.meedamian.bc1toolkit"
	DefaultConfig = "./bc1toolkit.conf"
)

type (
	Dialers struct {
		Tor      proxy.Dialer
		ClearNet proxy.Dialer
		mode     string
	}

	logger struct {
		name  string
		level logrus.Level
	}
)

func (l *logger) Get() *logrus.Entry {
	e := logrus.New()
	e.SetLevel(l.level)
	name := "unknown"
	if l.name != "" {
		name = l.name
	}

	return e.WithField("binary", name)
}

func (l *logger) Name(name string) {
	l.name = name
}

func (l *logger) SetLevel(level logrus.Level) {
	l.level = level
}

var Logger logger

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

		Logger.Get().WithError(err).Debugln("--tor-mode=auto and not Tor connection available: falling back to Clearnet only")
	}

	d.Tor = torDialer
	return
}

func GetCacheDir() string {
	return path.Join(cacheBase, cacheDir)
}
