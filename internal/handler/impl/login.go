package impl

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

func (i *implementation) Login(w http.ResponseWriter, r *http.Request) error {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		return err
	}

	var req RegisterReq
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

	sessionToken, err := i.svc.Login(ctx, *req.Login, *req.Password)
	if err != nil {
		return err
	}

	if sessionToken == nil {
		status = http.StatusUnauthorized
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   *sessionToken,
			Expires: time.Now().Add(24 * time.Hour),
		})
	}

	w.WriteHeader(status)
	return nil
}
