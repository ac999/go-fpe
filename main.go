// main.go
package main

import (
	"fmt"

	"github.com/ac999/go-fpe/algorithms"
)

func main() {
	K := []byte("examplekey123456") // 16 bytes for AES-128
	T := []byte("tweak")
	X := "1234567890"
	radix := 10
	minlen := 2
	maxlen := 32
	maxTlen := 16

	encrypted := algorithms.FF1Encrypt(K, T, X, radix, minlen, maxlen, maxTlen)
	fmt.Println("Encrypted:", encrypted)

	decrypted := algorithms.FF1Decrypt(K, T, encrypted, radix, minlen, maxlen, maxTlen)
	fmt.Println("Decrypted:", decrypted)
}
