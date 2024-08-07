package impl

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/chernyshevuser/gopfermart/internal/business"
)

func (a *api) NewOrder(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	status := http.StatusOK

	token, err := r.Cookie("token")
	if err != nil {
		status = http.StatusUnauthorized
		w.WriteHeader(status)
		return nil
	}

	if r.Header.Get("Content-Type") != "text/plain" {
		status = http.StatusBadRequest
		w.WriteHeader(status)
		return nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		status = http.StatusBadRequest
		w.WriteHeader(status)
		return nil
	}

	orderNum := strings.TrimSpace(string(body))

	err = a.svc.NewOrder(ctx, token.Value, orderNum)
	if err != nil {
		if errors.Is(err, business.ErrUnauthorized) {
			status = http.StatusUnauthorized
			w.WriteHeader(status)
			return nil
		}

		if errors.Is(err, business.ErrOrderRegisteredByUser) {
			w.WriteHeader(status)
			return nil
		}

		if errors.Is(err, business.ErrOrderRegisteredByOtherUser) {
			status = http.StatusConflict
			w.WriteHeader(status)
			return nil
		}

		if errors.Is(err, business.ErrIncorrectOrderNumber) {
			status = http.StatusUnprocessableEntity
			w.WriteHeader(status)
			return nil
		}
	}

	status = http.StatusAccepted
	w.WriteHeader(status)
	return nil
}
