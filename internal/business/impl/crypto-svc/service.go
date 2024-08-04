package crypto

type Svc interface {
	Encrypt(in string) (out string, err error)
	Decrypt(in string) (out string, err error)
}
