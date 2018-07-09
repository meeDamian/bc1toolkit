package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/meeDamian/bc1toolkit/lib/common"
)

var (
	Opts struct {
		Blab string `short:"b" long:"blab" description:"Blab blab blab"`
	}

	theirAddr = "192.168.1.221:18333"
)

func init() {
	common.Parser.AddGroup("bc1isup", "", &Opts)
	common.Parser.Parse()
	common.DefaultActions()
}

const (
	MainNet          uint32 = 0xd9b4bef9
	TestNet          uint32 = 0x0709110b

	ProtocolVersion  uint32 = 70013
	MaxPayloadLength uint32 = 358 // Don't ask why xD
	WitnessEncoding  uint32 = 2

	CommandSize    = 12
	VersionCommand = "version"
	UserAgent      = "/bc1isup:0.0.1/"
)

func main() {
	var bw bytes.Buffer

	// 4 bytes ; protocol version
	buf := make([]byte, 4)[:4]
	binary.LittleEndian.PutUint32(buf, ProtocolVersion)
	bw.Write(buf)

	// 8 bytes ; services enabled
	buf = make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, 0)
	bw.Write(buf)

	// 8 bytes ; timestamp
	buf = make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(time.Now().Unix()))
	bw.Write(buf)

	// 8 bytes ; their services
	buf = make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, 0)
	bw.Write(buf)

	// TODO: support v6…
	addr := strings.Split(theirAddr, ":")
	var port string
	ip := net.ParseIP(addr[0])
	if len(addr) > 1 {
		port = addr[1]
	}

	buf = make([]byte, 16)
	copy(buf[:], ip.To16())
	bw.Write(buf)

	// 2 bytes ; their port
	buf = make([]byte, 2)
	uintPort, _ := strconv.ParseUint(port, 10, 32)
	binary.BigEndian.PutUint16(buf, uint16(uintPort))
	bw.Write(buf)

	// 8 bytes ; our services
	buf = make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, 0)
	bw.Write(buf)

	// 16 bytes ; our IP address
	buf = make([]byte, 16)
	copy(buf[:], net.IPv6loopback.To16())
	bw.Write(buf)

	// 2 bytes ; our port
	buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(8333))
	bw.Write(buf)

	// 8 bytes ; nonce
	buf = make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(rand.Int63()))
	bw.Write(buf)

	// 1 byte ; user agent length
	buf = make([]byte, 1)
	buf[0] = uint8(len(UserAgent))
	bw.Write(buf)

	// ? bytes ; user agent string
	bw.Write([]byte(UserAgent))

	// 4 bytes ; last known block
	buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, 0)
	bw.Write(buf)

	// 1 byte ; disable tx relay
	buf = make([]byte, 1)
	buf[0] = uint8(0)
	bw.Write(buf)

	payload := bw.Bytes()

	// header
	hw := bytes.NewBuffer(make([]byte, 0, 24))

	// 4 bytes ; network magic
	buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, MainNet) // TODO: check both
	hw.Write(buf)

	// 12 bytes ; command
	var command [CommandSize]byte
	copy(command[:], VersionCommand)
	hw.Write(command[:])

	// 4 bytes ; payload length
	payloadLen := len(payload)

	buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(payloadLen))
	hw.Write(buf)

	// 4 bytes ; checksum
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	checksum := second[0:4]
	hw.Write(checksum)

	conn, err := net.Dial("tcp", theirAddr)
	if err != nil {
		fmt.Println("Couldn't connect to peer", err)
		return
	}

	fmt.Println("connected")

	// Write header.
	_, err = conn.Write(hw.Bytes())
	if err != nil {
		fmt.Println("couldn't .Write() header", err)
		return
	}

	fmt.Println("header sent")

	// write payload
	_, err = conn.Write(payload)
	if err != nil {
		fmt.Println("couldn't .Write() payload", err)
		return
	}

	fmt.Println("payload sent")

	reply := make([]byte, 100)
	_, err = conn.Read(reply)
	if err != nil {
		fmt.Println("xxx", err)
	}

	fmt.Println("reply from server=\n", string(reply))

	fmt.Println(hex.Dump(reply))

	conn.Close()
}
