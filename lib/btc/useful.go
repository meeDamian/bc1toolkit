package btc

import "crypto/sha256"

func DoubleSha256(v []byte) []byte {
	first := sha256.Sum256(v)
	second := sha256.Sum256(first[:])
	return second[:]
}
