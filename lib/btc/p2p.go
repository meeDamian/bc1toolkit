package btc

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/meeDamian/bc1toolkit/lib/common"
	"github.com/meeDamian/bc1toolkit/lib/connstring"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
)

const (
	MainNet uint32 = 0xd9b4bef9
	TestNet uint32 = 0x0709110b

	ProtocolVersion  uint32 = 70013
	MaxPayloadLength uint32 = 358 // Don't ask why ^^
	WitnessEncoding  uint32 = 2

	CommandSize    = 12
	HeaderSize     = 24
	VersionCommand = "version"
	VerAckCommand  = "verack"
	UserAgent      = "/bc1isup:0.0.1/"
)

type BitcoinVersion struct {
	Address   string `json:"address"`
	UserAgent string `json:"useragent"`
	Version   int    `json:"protocol"`
	LastBlock int    `json:"lastblock"`
	TestNet   bool   `json:"testnet"`
}

var defaultTimeout = time.Duration(5 * time.Second)

func getNetworkMagic(testNet bool) uint32 {
	if testNet {
		return TestNet
	}

	return MainNet
}

func buildNodeAddr(services uint64, ip net.IP, port string) []byte {
	var bw bytes.Buffer

	// 8 bytes ; services
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, services)
	bw.Write(buf)

	// 16 bytes ; IP address
	buf = make([]byte, 16)
	copy(buf[:], ip.To16())
	bw.Write(buf)

	// 2 bytes ; port
	buf = make([]byte, 2)
	uintPort, _ := strconv.ParseUint(port, 10, 32)
	binary.BigEndian.PutUint16(buf, uint16(uintPort))
	bw.Write(buf)

	return bw.Bytes()
}

func readNodeAddr(b []byte) (services uint64, ip net.IP, port uint16) {
	services = binary.LittleEndian.Uint64(b[:8])
	b = b[8:]

	ip = b[:16]
	b = b[16:]

	port = binary.BigEndian.Uint16(b[:2])

	return
}

func buildVersionHeader(msg []byte, testNet bool) []byte {
	magic := getNetworkMagic(testNet)

	checksum := DoubleSha256(msg)[0:4]

	common.Logger.Get().WithFields(logrus.Fields{
		"magic":    fmt.Sprintf("%x", magic),
		"cmd":      "version",
		"len":      len(msg),
		"checksum": checksum,
	}).Debugf("building header…")

	b := bytes.NewBuffer(make([]byte, 0, 24))

	// 4 bytes ; network magic
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, magic)
	b.Write(buf)

	// 12 bytes ; command
	var command [CommandSize]byte
	copy(command[:], VersionCommand)
	b.Write(command[:])

	// 4 bytes ; payload length
	payloadLen := len(msg)

	buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(payloadLen))
	b.Write(buf)

	// 4 bytes ; checksum
	b.Write(checksum)

	return b.Bytes()
}

func readHeader(header [HeaderSize]byte) (magic uint32, command string, length uint32, checksum [4]byte) {
	magic = binary.LittleEndian.Uint32(header[:4])
	header2 := header[4:]

	command = string(bytes.TrimRight(header2[:CommandSize], string(0)))
	header2 = header2[CommandSize:]

	length = binary.LittleEndian.Uint32(header2[:4])
	header2 = header2[4:]

	copy(checksum[:], header2[:4])
	return
}

func buildVersionMsg(ip net.IP, port string) []byte {
	var b bytes.Buffer

	// 4 bytes ; protocol version
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, ProtocolVersion)
	b.Write(buf)

	// 8 bytes ; services enabled
	buf = make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, 0)
	b.Write(buf)

	// 8 bytes ; timestamp
	buf = make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(time.Now().Unix()))
	b.Write(buf)

	// 26 bytes ; their address
	b.Write(buildNodeAddr(0, ip, port))

	// 26 bytes ; our address
	b.Write(buildNodeAddr(0, net.IPv6loopback, "0"))

	// 8 bytes ; nonce
	buf = make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(rand.Int63()))
	b.Write(buf)

	// 1 byte ; user agent length
	buf = make([]byte, 1)
	buf[0] = uint8(len(UserAgent))
	b.Write(buf)

	// len(UserAgent) bytes ; user agent string
	b.Write([]byte(UserAgent))

	// 4 bytes ; last known block
	buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, 0)
	b.Write(buf)

	// 1 byte ; disable tx relay
	buf = make([]byte, 1)
	buf[0] = uint8(0)
	b.Write(buf)

	return b.Bytes()
}

func readVersionMsg(msg []byte) (btcVersion BitcoinVersion) {
	btcVersion.Version = int(binary.LittleEndian.Uint32(msg[:4]))
	msg = msg[4:]

	//services = binary.LittleEndian.Uint64(msg[:8])
	msg = msg[8:]

	//ts = time.Time(time.Unix(int64(binary.LittleEndian.Uint64(msg[:8])), 0))
	msg = msg[8:]

	// our address
	//ourServices, ourIp, ourPort := readNodeAddr(msg[:26])
	msg = msg[26:]

	// their
	//theirServices, theirIp, theirPort := readNodeAddr(msg[:26])
	msg = msg[26:]

	// nonce (discard)
	msg = msg[8:]

	// user agent
	length := uint8(msg[0])
	if length == 0xff || length == 0xfe || length == 0xfd {
		// TODO: fixme
		fmt.Errorf("long (%d) useragents not yet supported…", length)
		return
	}
	msg = msg[1:]

	btcVersion.UserAgent = string(msg[:length])
	msg = msg[length:]

	if len(msg) > 0 {
		btcVersion.LastBlock = int(binary.LittleEndian.Uint32(msg[:4]))
	}

	return
}

func sendMessage(dialer proxy.Dialer, addr connstring.ConnString, header, msg []byte, testNet bool) (btcVersion BitcoinVersion, err error) {
	log := common.Logger.Get().WithField("address", addr.Raw)

	log.Debugln("connecting…")
	conn, err := dialer.Dial("tcp", net.JoinHostPort(addr.Host, addr.Port))
	if err != nil {
		return btcVersion, errors.Wrap(err, "can't connect to peer")
	}
	log.Debugln("connection ok")

	conn.SetDeadline(time.Now().Add(defaultTimeout))
	log.Debugf("read/write deadline set to %v", defaultTimeout)

	defer conn.Close()

	// write header
	log.WithField("header", fmt.Sprintf("%02x", header)).Debugln("sending header…")
	_, err = conn.Write(header)
	if err != nil {
		return btcVersion, errors.Wrap(err, "can't send header")
	}
	log.Debugln("header sent")

	// write payload
	log.WithField("payload", fmt.Sprintf("%02x", msg)).Debugln("sending payload…")
	_, err = conn.Write(msg)
	if err != nil {
		return btcVersion, errors.Wrap(err, "can't send payload")
	}
	log.Debugln("payload sent")

	// read header
	var respHeader [HeaderSize]byte
	_, err = conn.Read(respHeader[:])
	if err != nil {
		return btcVersion, errors.Wrap(err, "can't read peer header")
	}
	log.WithField("header", fmt.Sprintf("%02x", respHeader)).Debugln("peer header received")

	magic, command, length, checksum := readHeader(respHeader)
	log.WithFields(logrus.Fields{
		"magic":    fmt.Sprintf("%x", magic),
		"cmd":      command,
		"len":      length,
		"checksum": checksum,
	}).Debugln("peer header processed")

	if magic != MainNet && magic != TestNet {
		return btcVersion, errors.Errorf("peer node responded with a non-Bitcoin network magic (%02x)", magic)
	}

	if magic != getNetworkMagic(testNet) {
		return btcVersion, errors.Errorf("peer node responded with a unexpected network magic (expected:%02x, returned:%02x)", getNetworkMagic(testNet), magic)
	}

	if command != VersionCommand {
		return btcVersion, errors.New("peer node failed to reply correctly")
	}

	if length > MaxPayloadLength {
		return btcVersion, errors.New("possibly a malicious peer node detected")
	}

	respMsg := make([]byte, length)
	_, err = conn.Read(respMsg[:])
	if err != nil {
		return btcVersion, errors.Wrap(err, "can't read peer payload")
	}
	log.WithField("payload", fmt.Sprintf("%02x", respMsg)).Debugln("peer payload received")

	first := sha256.Sum256(respMsg[:])
	second := sha256.Sum256(first[:])
	if !bytes.Equal(checksum[:], second[0:4]) {
		return btcVersion, errors.New("received payload checksum does not match")
	}

	nodeInfo := readVersionMsg(respMsg)
	log.WithFields(logrus.Fields{
		"version":   nodeInfo.Version,
		"useragent": nodeInfo.UserAgent,
		"lastblock": nodeInfo.LastBlock,
	}).Debugln("peer payload processed")

	return nodeInfo, nil
}

func Speak(dialer proxy.Dialer, addr connstring.ConnString, testNet bool) (interface{}, error) {
	if addr.Type == connstring.TypeTorV3 {
		return BitcoinVersion{}, errors.New("Bitcoin network doesn't support Tor v3 addresses yet")
	}

	if addr.Port == "" {
		addr.Port = "8333"

		if testNet {
			addr.Port = "1" + addr.Port
		}
	}

	msg := buildVersionMsg(addr.IP, addr.Port)
	header := buildVersionHeader(msg, testNet)

	v := make(chan BitcoinVersion, 1)
	e := make(chan error, 1)

	go func() {
		version, err := sendMessage(dialer, addr, header, msg, testNet)
		if err != nil {
			e <- err
			return
		}

		if version == (BitcoinVersion{}) {
			e <- errors.New("empty version returned")
			return
		}

		version.Address = addr.ToString()
		version.TestNet = testNet

		v <- version
	}()

	select {
	case version := <-v:
		return version, nil

	case err := <-e:
		return BitcoinVersion{}, err

	case <-time.After(2 * defaultTimeout):
		return BitcoinVersion{}, errors.New("timeout")
	}
}
