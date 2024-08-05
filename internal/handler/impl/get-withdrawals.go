package impl

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"time"

	"github.com/chernyshevuser/gopfermart/internal/business"
)

type Withdrawal struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

func (a *api) GetWithdrawals(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	status := http.StatusOK

	token, err := r.Cookie("token")
	if err != nil {
		status = http.StatusUnauthorized
		w.WriteHeader(status)
		return nil
	}

	withdrawals, err := a.svc.GetWithdrawals(ctx, token.Value)
	if err != nil {
		if errors.Is(err, business.ErrUnauthorized) {
			status = http.StatusUnauthorized
			w.WriteHeader(status)
			return nil
		}
		return err
	}

	if len(withdrawals) == 0 {
		status = http.StatusNoContent
		w.WriteHeader(status)
		return nil
	}

	sort.Slice(withdrawals, func(i, j int) bool {
		return withdrawals[i].ProcessedAt.Before(withdrawals[j].ProcessedAt)
	})

	data := func() (out []Withdrawal) {
		for _, tmp := range withdrawals {
			out = append(out, Withdrawal{
				Order:       tmp.Order,
				Sum:         tmp.Sum,
				ProcessedAt: tmp.ProcessedAt.Format(time.RFC3339),
			})
		}
		return out
	}()

	resp, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	return nil
}
