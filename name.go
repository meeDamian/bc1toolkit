package main

import (
	"fmt"
	"strings"

	"github.com/btcsuite/btcutil/bech32"
)

const (
	charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"
	hrp     = "bc"
	name    = "tools"
)

func main() {
	var word []byte
	for _, letter := range name {
		l := string(letter)

		if l == "o" {
			l = "0"
		}

		idx := strings.Index(charset, l)
		if idx == -1 {
			panic(fmt.Sprintf("%s is not allowed in Bech32 charset", l))
		}

		word = append(word, byte(idx))
	}

	// TODO: replace with own
	encoded, err := bech32.Encode(hrp, word)
	if err != nil {
		panic(err)
	}

	fmt.Println(encoded)
}
