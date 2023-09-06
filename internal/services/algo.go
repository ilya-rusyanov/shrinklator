package services

import (
	"crypto/md5"
	"encoding/hex"
)

type Algo = func(string) string

func MD5Algo(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
