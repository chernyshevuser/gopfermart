package sessionsvc

// Svc handles with access tokens
type Svc interface {
	NewToken(login string) (token string, err error)
	GetLogin(token string) (login string, ok bool)
}
