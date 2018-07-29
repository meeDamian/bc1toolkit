package tor

import (
	"encoding/json"
		"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/proxy"
)

func checkConnection(d proxy.Dialer) error {
	client := http.Client{
		Transport: &http.Transport{
			Dial: d.Dial,
		},
		Timeout: 4 * time.Second,
	}

	// TODO: remove the need for this request(?)
	resp, err := client.Get("https://check.torproject.org/api/ip")
	if err != nil {
		return errors.Wrap(err, "can't .Get() Tor-check website")
	}

	v := struct{ IsTor bool `json:"IsTor"` }{}
	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		return errors.Wrap(err, "Unable to JSON-decode response")
	}

	if !v.IsTor {
		return errors.New("please tell me when this happens…")
	}

	return nil
}

func GetWorkingTor(addresses []string) (dialer proxy.Dialer, err error) {
	for _, addr := range addresses {
		dialer, err = proxy.SOCKS5("tcp", addr, nil, nil)
		if err != nil {
			continue
		}

		err = checkConnection(dialer)
		if err == nil {
			return dialer, nil
		}
	}

	return nil, errors.Wrap(err, "can't find a working Tor connection")
}
