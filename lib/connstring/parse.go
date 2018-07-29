package connstring

import (
	"fmt"
	"net"
	"strings"

	"github.com/pkg/errors"
)

const (
	_ = iota
	TypeIpV4
	TypeIpV6
	TypeTorV2
	TypeTorV3
	TypeDomain

	onionTld  = ".onion"
	localTld  = ".local"
	localhost = "localhost"

	pubKeyLen = 66
)

type ConnString struct {
	Raw    string
	PubKey string
	Addr   string
	IP     net.IP
	Port   string
	Type   int
	Local  bool
}

var localMasks = []string{
	"127.0.0.0/8",    // IPv4 loopback
	"10.0.0.0/8",     // RFC1918
	"172.16.0.0/12",  // RFC1918
	"192.168.0.0/16", // RFC1918
	"::1/128",        // IPv6 loopback
	"fe80::/10",      // IPv6 link-local
}

func (c ConnString) IsTor() bool {
	return c.Type == TypeTorV2 || c.Type == TypeTorV3
}

func (c ConnString) ToString() (connstring string) {
	if c.PubKey != "" {
		connstring = fmt.Sprintf("%s@", c.PubKey)
	}

	connstring += c.Addr

	if c.Port != "" {
		connstring += fmt.Sprintf(":%s", c.Port)
	}

	return
}

func Parse(connstring string) (c ConnString, err error) {
	c.Raw = connstring

	connstring = strings.ToLower(connstring)
	connstring = strings.TrimPrefix(connstring, "http://")
	connstring = strings.TrimPrefix(connstring, "https://")

	// Process LN pubkey if available
	chunks := strings.Split(connstring, "@")
	if len(chunks) > 1 {
		if len(chunks[0]) != pubKeyLen {
			return c, errors.Errorf("Invalid LN pubkey length: %d instead of %d", len(chunks[0]), pubKeyLen)
		}

		c.PubKey = chunks[0]
		connstring = chunks[1]

	} else {
		connstring = chunks[0]
	}

	// split address into ip/host and port
	c.Addr = connstring
	if (strings.Count(connstring, ":") > 1 && strings.Contains(connstring, "[")) || strings.Count(connstring, ":") == 1 {
		c.Addr, c.Port, err = net.SplitHostPort(connstring)
		if err != nil {
			return
		}
	}

	// process Tor
	if strings.Contains(c.Addr, onionTld) {
		c.Type = TypeTorV2

		if len(c.Addr) == 56+len(onionTld) {
			c.Type = TypeTorV3
		}

		return
	}

	// process IPs
	c.IP = net.ParseIP(c.Addr)
	if c.IP != nil {
		c.Type = TypeIpV4

		if c.IP.To4() == nil {
			c.Type = TypeIpV6
		}

		// check if given connstring is a local one
		for _, cidr := range localMasks {
			_, block, _ := net.ParseCIDR(cidr)

			if block.Contains(c.IP) {
				c.Local = true
			}
		}

		return
	}

	// process domains
	c.Type = TypeDomain
	if c.Addr == localhost || strings.Contains(c.Addr, localTld) {
		c.Local = true
	}

	return
}
