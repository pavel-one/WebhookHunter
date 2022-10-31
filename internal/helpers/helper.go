package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TrimJson(jsonBytes []byte) []byte {
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, jsonBytes); err != nil {
		fmt.Println(err)
	}

	return buffer.Bytes()
}
