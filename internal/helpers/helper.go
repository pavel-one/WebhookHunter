package helpers

import (
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz123456789"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
