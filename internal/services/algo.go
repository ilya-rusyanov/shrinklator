package services

import (
	"crypto/md5"
	"encoding/hex"
)

// Algo - prototype for shortening algorithm
type Algo = func(string) string

// MD5Algo md5 implementation of shortening algorithm
func MD5Algo(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
