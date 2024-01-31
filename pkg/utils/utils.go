package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"github.com/xdg-go/pbkdf2"
)

var (
	Uppercase        = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Lowercase        = "abcdefghijklmnopqrstuvwxyz"
	Numeric          = "0123456789"
	UppercaseNumeric = fmt.Sprintf("%s%s", Uppercase, Numeric)
	Alpha            = fmt.Sprintf("%s%s", Uppercase, Lowercase)
	AlphaNumeric     = fmt.Sprintf("%s%s", Alpha, Numeric)
)

func GenerateUniqueID(composition string, length int) (id string) {
	if len(composition) < 1 {
		return
	}
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("crypto/rand.Read error: %v", err))
	}
	for i, byt := range b {
		b[i] = composition[int(byt)%len(composition)]
	}
	return string(b)
}

func Encrypt(plainText, salt string) string {
	cryptedBuf := pbkdf2.Key([]byte(plainText), []byte(salt), 32, 64, sha512.New)
	crypted := hex.EncodeToString(cryptedBuf)

	return crypted
}
