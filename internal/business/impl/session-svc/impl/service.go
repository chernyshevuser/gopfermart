package impl

import (
	"sync"
	"time"

	sessionsvc "github.com/chernyshevuser/gopfermart.git/internal/business/impl/session-svc"
	"github.com/golang-jwt/jwt"
)

type service struct {
	jwtKey  []byte
	storage map[string]string
	// assume read locks occur more frequently than write locks
	mu *sync.RWMutex
}

func New(secretKey string) sessionsvc.Svc {
	return &service{
		jwtKey:  []byte(secretKey),
		storage: make(map[string]string),
		mu:      &sync.RWMutex{},
	}
}

func (s *service) NewToken(login string) (token string, err error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Login: login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = jwtToken.SignedString(s.jwtKey)
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.storage[login] = token

	return token, nil
}

func (s *service) CheckToken(login string, token string) (ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	storedToken, exists := s.storage[login]
	if !exists {
		return false
	}

	return token == storedToken
}
