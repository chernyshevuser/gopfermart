package utils

import (
	"github.com/chernyshevuser/gopfermart.git/tools/config"
	"github.com/chernyshevuser/gopfermart.git/tools/crypto"
)

func Encrypt(in string) (out string, err error) {
	return crypto.Encrypt(config.CryptoKey, in)
}

func Decrypt(in string) (out string, err error) {
	return crypto.Decrypt(config.CryptoKey, in)
}
