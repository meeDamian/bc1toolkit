package connstring

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	pubkey = "032260c3b64b471b7eb0630b4af5d07ca94ff4e759573cbbe1bfb25845c375ed6e"
	port   = "8333"

	ipV4NoPort   = "112.108.141.146"
	ipV4WithPort = ipV4NoPort + ":" + port

	localIpV4NoPort   = "192.168.1.100"
	localIpV4WithPort = localIpV4NoPort + ":" + port

	ipV6NoPort   = "8323:55b0:587e:899c:ecbb:835b:ef35:c67"
	ipV6WithPort = "[" + ipV6NoPort + "]:" + port

	localIpV6NoPort   = "fe80::cab:caa5:d1f2:6947"
	localIpV6WithPort = "[" + localIpV6NoPort + "]:" + port

	torV2NoPort   = "s7eo5ftart6m3tmm.onion"
	torV2WithPort = torV2NoPort + ":" + port

	torV3NoPort     = "6g3y7ahr5uxjzedgmu5etxtrc6hqmb2bl7kyipl3nzrcnyp64afrfwyd.onion"
	torV3WithPort   = torV3NoPort + ":" + port
	lnTorV3WithPort = pubkey + "@" + torV3WithPort

	domainNoPort   = "example.com"
	domainWithPort = "https://" + domainNoPort + ":" + port

	localDomainNoPort   = "example.local"
	localDomainWithPort = localDomainNoPort + ":" + port

	localhostNoPort   = "localhost"
	localhostWithPort = localhostNoPort + ":" + port
)

// TODO: add error cases

//
// IP v4
//
func TestParseIpV4(t *testing.T) {
	Convey("Given a valid IPv4 address w/o a port", t, func() {
		addr, err := Parse(ipV4NoPort)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("addr should not be empty", func() {
			So(addr, ShouldNotBeEmpty)
		})

		Convey(".Type should be address.TypeIpV4", func() {
			So(addr.Type, ShouldEqual, TypeIpV4)
		})

		Convey(".Addr should be set to the original IP", func() {
			So(addr.Addr, ShouldEqual, ipV4NoPort)
		})

		Convey(".Port should be empty", func() {
			So(addr.Port, ShouldBeEmpty)
		})

		Convey(".IP should be set correctly", func() {
			So(addr.IP.String(), ShouldEqual, ipV4NoPort)
		})

		Convey(".Local should be false", func() {
			So(addr.Local, ShouldBeFalse)
		})
	})

	Convey("Given a valid IPv4 address with a port", t, func() {
		addr, err := Parse(ipV4WithPort)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("addr should not be empty", func() {
			So(addr, ShouldNotBeEmpty)
		})

		Convey(".Type should be address.TypeIpV4", func() {
			So(addr.Type, ShouldEqual, TypeIpV4)
		})

		Convey(".Addr should be set to the original IP", func() {
			So(addr.Addr, ShouldEqual, ipV4NoPort)
		})

		Convey(".Port should be set to 8333", func() {
			So(addr.Port, ShouldEqual, port)
		})

		Convey(".IP should be set correctly", func() {
			So(addr.IP.String(), ShouldEqual, ipV4NoPort)
		})

		Convey(".Local should be false", func() {
			So(addr.Local, ShouldBeFalse)
		})
	})
}

func TestParseLocalIpV4(t *testing.T) {
	Convey("Given a valid IPv4 address w/o a port", t, func() {
		addr, err := Parse(localIpV4NoPort)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("addr should not be empty", func() {
			So(addr, ShouldNotBeEmpty)
		})

		Convey(".Type should be address.TypeIpV4", func() {
			So(addr.Type, ShouldEqual, TypeIpV4)
		})

		Convey(".Addr should be set to the original IP", func() {
			So(addr.Addr, ShouldEqual, localIpV4NoPort)
		})

		Convey(".Port should be empty", func() {
			So(addr.Port, ShouldBeEmpty)
		})

		Convey(".IP should be set correctly", func() {
			So(addr.IP.String(), ShouldEqual, localIpV4NoPort)
		})

		Convey(".Local should be true", func() {
			So(addr.Local, ShouldBeTrue)
		})
	})

	Convey("Given a valid IPv4 address with a port", t, func() {
		addr, err := Parse(localIpV4WithPort)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("addr should not be empty", func() {
			So(addr, ShouldNotBeEmpty)
		})

		Convey(".Type should be address.TypeIpV4", func() {
			So(addr.Type, ShouldEqual, TypeIpV4)
		})

		Convey(".Addr should be set to the original IP", func() {
			So(addr.Addr, ShouldEqual, localIpV4NoPort)
		})

		Convey(".Port should be set to 8333", func() {
			So(addr.Port, ShouldEqual, port)
		})

		Convey(".IP should be set correctly", func() {
			So(addr.IP.String(), ShouldEqual, localIpV4NoPort)
		})

		Convey(".Local should be false", func() {
			So(addr.Local, ShouldBeTrue)
		})
	})
}

//
// IP v6
//
func TestParseIpV6(t *testing.T) {
	Convey("Given a valid IPv6 address w/o a port", t, func() {
		addr, err := Parse(ipV6NoPort)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("addr should not be empty", func() {
			So(addr, ShouldNotBeEmpty)
		})

		Convey(".Type should be address.TypeIpV6", func() {
			So(addr.Type, ShouldEqual, TypeIpV6)
		})

		Convey(".Addr should be set to the original IP", func() {
			So(addr.Addr, ShouldEqual, ipV6NoPort)
		})

		Convey(".Port should be empty", func() {
			So(addr.Port, ShouldBeEmpty)
		})

		Convey(".IP should be set correctly", func() {
			So(addr.IP.String(), ShouldEqual, ipV6NoPort)
		})

		Convey(".Local should be false", func() {
			So(addr.Local, ShouldBeFalse)
		})
	})

	Convey("Given a valid IPv6 address with a port", t, func() {
		addr, err := Parse(ipV6WithPort)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("addr should not be empty", func() {
			So(addr, ShouldNotBeEmpty)
		})

		Convey(".Type should be address.TypeIpV6", func() {
			So(addr.Type, ShouldEqual, TypeIpV6)
		})

		Convey(".Addr should be set to the original IP", func() {
			So(addr.Addr, ShouldEqual, ipV6NoPort)
		})

		Convey(".Port should be set to 8333", func() {
			So(addr.Port, ShouldEqual, port)
		})

		Convey(".IP should be set correctly", func() {
			So(addr.IP.String(), ShouldEqual, ipV6NoPort)
		})

		Convey(".Local should be false", func() {
			So(addr.Local, ShouldBeFalse)
		})
	})
}

func TestParseLocalIpV6(t *testing.T) {
	Convey("Given a valid IPv6 address w/o a port", t, func() {
		addr, err := Parse(localIpV6NoPort)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("addr should not be empty", func() {
			So(addr, ShouldNotBeEmpty)
		})

		Convey(".Type should be address.TypeIpV6", func() {
			So(addr.Type, ShouldEqual, TypeIpV6)
		})

		Convey(".Addr should be set to the original IP", func() {
			So(addr.Addr, ShouldEqual, localIpV6NoPort)
		})

		Convey(".Port should be empty", func() {
			So(addr.Port, ShouldBeEmpty)
		})

		Convey(".IP should be set correctly", func() {
			So(addr.IP.String(), ShouldEqual, localIpV6NoPort)
		})

		Convey(".Local should be true", func() {
			So(addr.Local, ShouldBeTrue)
		})
	})

	Convey("Given a valid IPv6 address with a port", t, func() {
		addr, err := Parse(localIpV6WithPort)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("addr should not be empty", func() {
			So(addr, ShouldNotBeEmpty)
		})

		Convey(".Type should be address.TypeIpV6", func() {
			So(addr.Type, ShouldEqual, TypeIpV6)
		})

		Convey(".Addr should be set to the original IP", func() {
			So(addr.Addr, ShouldEqual, localIpV6NoPort)
		})

		Convey(".Port should be set to 8333", func() {
			So(addr.Port, ShouldEqual, port)
		})

		Convey(".IP should be set correctly", func() {
			So(addr.IP.String(), ShouldEqual, localIpV6NoPort)
		})

		Convey(".Local should be false", func() {
			So(addr.Local, ShouldBeTrue)
		})
	})
}

//
// Tor v2
//
func TestParseTorV2(t *testing.T) {
	Convey("Given a valid TorV2 address w/o a port", t, func() {
		addr, err := Parse(torV2NoPort)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("addr should not be empty", func() {
			So(addr, ShouldNotBeEmpty)
		})

		Convey(".Type should be address.TypeTorV2", func() {
			So(addr.Type, ShouldEqual, TypeTorV2)
		})

		Convey(".Addr should be set to the original address", func() {
			So(addr.Addr, ShouldEqual, torV2NoPort)
		})

		Convey(".Port should be empty", func() {
			So(addr.Port, ShouldBeEmpty)
		})

		Convey(".IP should be empty", func() {
			So(addr.IP, ShouldBeEmpty)
		})

		Convey(".Local should be false", func() {
			So(addr.Local, ShouldBeFalse)
		})
	})

	Convey("Given a valid TorV2 address with a port", t, func() {
		addr, err := Parse(torV2WithPort)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("addr should not be empty", func() {
			So(addr, ShouldNotBeEmpty)
		})

		Convey(".Type should be address.TypeTorV2", func() {
			So(addr.Type, ShouldEqual, TypeTorV2)
		})

		Convey(".Addr should be set to the original address", func() {
			So(addr.Addr, ShouldEqual, torV2NoPort)
		})

		Convey(".Port should be set to 8333", func() {
			So(addr.Port, ShouldEqual, port)
		})

		Convey(".IP should be empty", func() {
			So(addr.IP, ShouldBeEmpty)
		})

		Convey(".Local should be false", func() {
			So(addr.Local, ShouldBeFalse)
		})
	})
}

//
// Tor v3
//
func TestParseTorV3(t *testing.T) {
	Convey("Given a valid TorV3 address w/o a port", t, func() {
		addr, err := Parse(torV3NoPort)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("addr should not be empty", func() {
			So(addr, ShouldNotBeEmpty)
		})

		Convey(".Type should be address.TypeTorV3", func() {
			So(addr.Type, ShouldEqual, TypeTorV3)
		})

		Convey(".Addr should be set to the original address", func() {
			So(addr.Addr, ShouldEqual, torV3NoPort)
		})

		Convey(".Port should be empty", func() {
			So(addr.Port, ShouldBeEmpty)
		})

		Convey(".IP should be empty", func() {
			So(addr.IP, ShouldBeEmpty)
		})

		Convey(".Local should be false", func() {
			So(addr.Local, ShouldBeFalse)
		})
	})

	Convey("Given a valid TorV3 address with a port", t, func() {
		addr, err := Parse(torV3WithPort)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("addr should not be empty", func() {
			So(addr, ShouldNotBeEmpty)
		})

		Convey(".Type should be address.TypeTorV3", func() {
			So(addr.Type, ShouldEqual, TypeTorV3)
		})

		Convey(".Addr should be set to the original IP", func() {
			So(addr.Addr, ShouldEqual, torV3NoPort)
		})

		Convey(".Port should be set to 8333", func() {
			So(addr.Port, ShouldEqual, port)
		})

		Convey(".IP should be empty", func() {
			So(addr.IP, ShouldBeEmpty)
		})

		Convey(".Local should be false", func() {
			So(addr.Local, ShouldBeFalse)
		})
	})

	Convey("Given a valid LN + TorV3 address with a port", t, func() {
		addr, err := Parse(lnTorV3WithPort)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("addr should not be empty", func() {
			So(addr, ShouldNotBeEmpty)
		})

		Convey(".PubKey should be set properly", func() {
			So(addr.PubKey, ShouldEqual, pubkey)
		})

		Convey(".Type should be address.TypeTorV3", func() {
			So(addr.Type, ShouldEqual, TypeTorV3)
		})

		Convey(".Addr should be set to the original IP", func() {
			So(addr.Addr, ShouldEqual, torV3NoPort)
		})

		Convey(".Port should be set to 8333", func() {
			So(addr.Port, ShouldEqual, port)
		})

		Convey(".IP should be empty", func() {
			So(addr.IP, ShouldBeEmpty)
		})

		Convey(".Local should be false", func() {
			So(addr.Local, ShouldBeFalse)
		})
	})
}

//
// Domain
//
func TestParseDomain(t *testing.T) {
	Convey("Given a valid domain w/o a port", t, func() {
		addr, err := Parse(domainNoPort)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("addr should not be empty", func() {
			So(addr, ShouldNotBeEmpty)
		})

		Convey(".Type should be address.TypeDomain", func() {
			So(addr.Type, ShouldEqual, TypeDomain)
		})

		Convey(".Addr should be set to the original domain", func() {
			So(addr.Addr, ShouldEqual, domainNoPort)
		})

		Convey(".Port should be empty", func() {
			So(addr.Port, ShouldBeEmpty)
		})

		Convey(".IP should be empty", func() {
			So(addr.IP, ShouldBeEmpty)
		})

		Convey(".Local should be false", func() {
			So(addr.Local, ShouldBeFalse)
		})
	})

	Convey("Given a valid domain with a port", t, func() {
		addr, err := Parse(domainWithPort)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("addr should not be empty", func() {
			So(addr, ShouldNotBeEmpty)
		})

		Convey(".Type should be address.TypeDomain", func() {
			So(addr.Type, ShouldEqual, TypeDomain)
		})

		Convey(".Addr should be set to the original domain", func() {
			So(addr.Addr, ShouldEqual, domainNoPort)
		})

		Convey(".Port should be set to 8333", func() {
			So(addr.Port, ShouldEqual, port)
		})

		Convey(".IP should be empty", func() {
			So(addr.IP, ShouldBeEmpty)
		})

		Convey(".Local should be false", func() {
			So(addr.Local, ShouldBeFalse)
		})
	})
}

//
// Local Domain
//
func TestParseLocalDomain(t *testing.T) {
	Convey("Given a valid local domain w/o a port", t, func() {
		addr, err := Parse(localDomainNoPort)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("addr should not be empty", func() {
			So(addr, ShouldNotBeEmpty)
		})

		Convey(".Type should be address.TypeDomain", func() {
			So(addr.Type, ShouldEqual, TypeDomain)
		})

		Convey(".Addr should be set to the original domain", func() {
			So(addr.Addr, ShouldEqual, localDomainNoPort)
		})

		Convey(".Port should be empty", func() {
			So(addr.Port, ShouldBeEmpty)
		})

		Convey(".IP should be empty", func() {
			So(addr.IP, ShouldBeEmpty)
		})

		Convey(".Local should be true", func() {
			So(addr.Local, ShouldBeTrue)
		})
	})

	Convey("Given a valid domain with a port", t, func() {
		addr, err := Parse(localDomainWithPort)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("addr should not be empty", func() {
			So(addr, ShouldNotBeEmpty)
		})

		Convey(".Type should be address.TypeDomain", func() {
			So(addr.Type, ShouldEqual, TypeDomain)
		})

		Convey(".Addr should be set to the original domain", func() {
			So(addr.Addr, ShouldEqual, localDomainNoPort)
		})

		Convey(".Port should be set to 8333", func() {
			So(addr.Port, ShouldEqual, port)
		})

		Convey(".IP should be empty", func() {
			So(addr.IP, ShouldBeEmpty)
		})

		Convey(".Local should be true", func() {
			So(addr.Local, ShouldBeTrue)
		})
	})
}

//
// localhost
//
func TestParseLocalhost(t *testing.T) {
	Convey("Given a valid localhost w/o a port", t, func() {
		addr, err := Parse(localhostNoPort)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("addr should not be empty", func() {
			So(addr, ShouldNotBeEmpty)
		})

		Convey(".Type should be address.TypeDomain", func() {
			So(addr.Type, ShouldEqual, TypeDomain)
		})

		Convey(".Addr should be set to the original domain", func() {
			So(addr.Addr, ShouldEqual, localhostNoPort)
		})

		Convey(".Port should be empty", func() {
			So(addr.Port, ShouldBeEmpty)
		})

		Convey(".IP should be empty", func() {
			So(addr.IP, ShouldBeEmpty)
		})

		Convey(".Local should be true", func() {
			So(addr.Local, ShouldBeTrue)
		})
	})

	Convey("Given a valid domain with a port", t, func() {
		addr, err := Parse(localhostWithPort)

		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})

		Convey("addr should not be empty", func() {
			So(addr, ShouldNotBeEmpty)
		})

		Convey(".Type should be address.TypeDomain", func() {
			So(addr.Type, ShouldEqual, TypeDomain)
		})

		Convey(".Addr should be set to the original domain", func() {
			So(addr.Addr, ShouldEqual, localhostNoPort)
		})

		Convey(".Port should be set to 8333", func() {
			So(addr.Port, ShouldEqual, port)
		})

		Convey(".IP should be empty", func() {
			So(addr.IP, ShouldBeEmpty)
		})

		Convey(".Local should be true", func() {
			So(addr.Local, ShouldBeTrue)
		})
	})
}
