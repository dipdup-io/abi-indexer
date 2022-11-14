package random

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"math/big"
)

var letters = []byte("-/0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

// String - returns random string with fixed length `n`
func String(n int) (string, error) {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

// UInt64 - returns random uint64 number
func UInt64() (uint64, error) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.ReadUvarint(bytes.NewBuffer(b))
}
