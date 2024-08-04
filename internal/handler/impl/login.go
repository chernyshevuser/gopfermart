package impl

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/chernyshevuser/gopfermart/internal/business"
)

type LoginReq struct {
	Login    *string `json:"login,omitempty"`
	Password *string `json:"password,omitempty"`
}

func (a *api) Login(w http.ResponseWriter, r *http.Request) error {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		return err
	}

	var req LoginReq
	if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
		return err
	}

	status := http.StatusOK

	if req.Login == nil || req.Password == nil {
		status = http.StatusBadRequest
		w.WriteHeader(status)
		return nil
	}

	ctx := r.Context()

	sessionToken, err := a.svc.Login(ctx, *req.Login, *req.Password)
	if err != nil {
		if errors.Is(err, business.ErrUnauthorized) {
			status = http.StatusUnauthorized
			w.WriteHeader(status)
			return nil
		}
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   sessionToken,
		Expires: time.Now().Add(24 * time.Hour),
	})

	w.WriteHeader(status)
	return nil
}
