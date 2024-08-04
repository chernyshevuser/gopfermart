package utils

import (
	cryptosvc "github.com/chernyshevuser/gopfermart/internal/business/impl/crypto-svc"
	"github.com/chernyshevuser/gopfermart/tools/config"
	"github.com/chernyshevuser/gopfermart/tools/crypto"
)

type svc struct {
	cryptoKey string
}

func New(cryptoKey string) cryptosvc.Svc {
	return &svc{
		cryptoKey: cryptoKey,
	}
}

func (s *svc) Encrypt(in string) (out string, err error) {
	return crypto.Encrypt(config.CryptoKey, in)
}

func (s *svc) Decrypt(in string) (out string, err error) {
	return crypto.Decrypt(config.CryptoKey, in)
}
