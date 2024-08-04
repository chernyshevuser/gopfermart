package impl

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/chernyshevuser/gopfermart/internal/business"
)

type GetBalanceResp struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func (a *api) GetBalance(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	status := http.StatusOK

	token, err := r.Cookie("token")
	if err != nil {
		status = http.StatusUnauthorized
		w.WriteHeader(status)
		return nil
	}

	current, withdrawn, err := a.svc.GetBalance(ctx, token.Value)
	if err != nil {
		if errors.Is(err, business.ErrUnauthorized) {
			status = http.StatusUnauthorized
			w.WriteHeader(status)
			return nil
		}
		return err
	}

	resp, err := json.Marshal(GetBalanceResp{
		Current:   current,
		Withdrawn: withdrawn,
	})
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	return nil
}
