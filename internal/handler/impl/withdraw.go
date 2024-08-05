package impl

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/chernyshevuser/gopfermart/internal/business"
)

type WithdrawReq struct {
	Order *string  `json:"order,omitempty"`
	Sum   *float64 `json:"sum,omitempty"`
}

func (a *api) Withdraw(w http.ResponseWriter, r *http.Request) error {
	status := http.StatusOK

	token, err := r.Cookie("token")
	if err != nil {
		status = http.StatusUnauthorized
		w.WriteHeader(status)
		return nil
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		return err
	}

	var req WithdrawReq
	if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
		return err
	}

	if req.Order == nil || req.Sum == nil {
		status = http.StatusBadRequest
		w.WriteHeader(status)
		return nil
	}

	ctx := r.Context()

	err = a.svc.Withdraw(ctx, token.Value, *req.Order, *req.Sum)
	if err != nil {
		if errors.Is(err, business.ErrUnauthorized) {
			status = http.StatusUnauthorized
			w.WriteHeader(status)
			return nil
		}
		if errors.Is(err, business.ErrInsufficientBalance) {
			status = http.StatusPaymentRequired
			w.WriteHeader(status)
			return nil
		}
		if errors.Is(err, business.ErrIncorrectOrderNumber) {
			status = http.StatusUnprocessableEntity
			w.WriteHeader(status)
			return nil
		}
		return err
	}

	w.WriteHeader(status)
	return nil
}
