package shortener

import (
	"errors"
	"math/big"
)

const base = 62

type Shortener struct {
	short2long map[int64]string
	long2short map[string]int64
}

func (s *Shortener) Shrink(input string) (string, error) {
	var hash int64

	if short, ok := s.long2short[input]; ok {
		hash = short
	} else {
		hash = int64(len(s.long2short))
		s.long2short[input] = hash
		s.short2long[hash] = input
	}

	var b62 big.Int
	b62.SetInt64(hash)
	return b62.Text(base), nil
}

var errUnknown = errors.New("unknown")
var errHashing = errors.New("bad input")

func (s *Shortener) Expand(input string) (string, error) {
	hasher := big.Int{}

	if _, ok := hasher.SetString(input, base); !ok {
		return "", errHashing
	}

	hash := hasher.Int64()

	if result, ok := s.short2long[hash]; ok {
		return result, nil
	}

	return "", errUnknown
}

func New() *Shortener {
	res := &Shortener{}
	res.short2long = map[int64]string{}
	res.long2short = map[string]int64{}
	return res
}
