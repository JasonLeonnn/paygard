package http

import (
	"encoding/json"
	nethttp "net/http"
	"time"

	"github.com/JasonLeonnn/jalytics/internal/db"
	"github.com/JasonLeonnn/jalytics/internal/services"
)

type CreateTransactionRequest struct {
	Amount          float64   `json:"amount"`
	Category        string    `json:"category"`
	Merchant        string    `json:"merchant"`
	TransactionDate time.Time `json:"transaction_date"`
}

func CreateTransactionHandler(service *services.TransactionService) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req CreateTransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			nethttp.Error(w, "Invalid request payload", nethttp.StatusBadRequest)
			return
		}

		if req.Amount <= 0 || req.Category == "" || req.TransactionDate.IsZero() {
			nethttp.Error(w, "Invalid transaction fields", nethttp.StatusBadRequest)
			return
		}

		tx := &db.Transaction{
			Amount:          req.Amount,
			Category:        req.Category,
			Merchant:        req.Merchant,
			TransactionDate: req.TransactionDate,
		}

		isAnomaly, id, err := service.CreateTransaction(r.Context(), tx)
		if err != nil {
			nethttp.Error(w, "Failed to create transaction", nethttp.StatusInternalServerError)
			return
		}

		w.WriteHeader(nethttp.StatusCreated)
		resp := map[string]interface{}{
			"id":      id,
			"anomaly": isAnomaly,
		}
		json.NewEncoder(w).Encode(resp)
	}
}
