package impl

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"time"

	"github.com/chernyshevuser/gopfermart/internal/business"
)

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    *int      `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func (a *api) GetOrders(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	status := http.StatusOK

	token, err := r.Cookie("token")
	if err != nil {
		status = http.StatusUnauthorized
		w.WriteHeader(status)
		return nil
	}

	orders, err := a.svc.GetOrders(ctx, token.Value)
	if err != nil {
		if errors.Is(err, business.ErrUnauthorized) {
			status = http.StatusUnauthorized
			w.WriteHeader(status)
			return nil
		}
		return err
	}

	if len(orders) == 0 {
		status = http.StatusNoContent
		w.WriteHeader(status)
		return nil
	}

	data, err := convertOrders(orders)
	if err != nil {
		return err
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].UploadedAt.Before(data[j].UploadedAt)
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

func convertOrders(in []business.Order) (out []Order, err error) {
	for _, order := range in {
		uploadedAt := order.UploadedAt
		if err != nil {
			return []Order{}, err
		}

		var accrual *int
		if order.Accrual != nil {
			accrualInt := int(*order.Accrual)
			accrual = &accrualInt
		}

		out = append(out, Order{
			Number:     order.Number,
			Status:     order.Status,
			Accrual:    accrual,
			UploadedAt: uploadedAt,
		})
	}

	return out, nil
}
