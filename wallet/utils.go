package wallet

import (
	"log"

	"github.com/mr-tron/base58"
)

// Base58Encode - это функция, которая кодирует в base58.
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)
	return []byte(encode)
}

// Base58Decode - это функция, которая декодирует из base58.
func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input))
	Handle(err)
	return decode
}

// Handle - это функция, которая обрабатывает ошибки.
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
