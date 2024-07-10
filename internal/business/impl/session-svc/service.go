package sessionsvc

// Svc handles with access tokens
type Svc interface {
	NewToken(login string) (token string, err error)
	CheckToken(login string, token string) (ok bool)
}
