package services

import "testing"

func BenchmarkMD5Algo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MD5Algo("sample")
	}
}
