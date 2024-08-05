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
	Order        string    `json:"order"`
	Sum          float64   `json:"sum"`
	Processed_at time.Time `json:"processed_at"`
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

	data, err := convertWithdrawals(withdrawals)
	if err != nil {
		return err
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Processed_at.Before(data[j].Processed_at)
	})

	resp, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	return nil
}

func convertWithdrawals(in []business.Withdrawal) (out []Withdrawal, err error) {
	for _, tmp := range in {
		processedAt, err := time.Parse(time.RFC3339, tmp.ProcessedAt)
		if err != nil {
			return []Withdrawal{}, err
		}

		out = append(out, Withdrawal{
			Order:        tmp.Order,
			Sum:          tmp.Sum,
			Processed_at: processedAt,
		})
	}

	return out, nil
}
