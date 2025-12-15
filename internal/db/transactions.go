package db

import (
	"context"
	"time"

	"github.com/JasonLeonnn/paygard/internal/metrics"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transaction struct {
	ID              string
	Amount          float64
	Category        string
	Merchant        string
	TransactionDate time.Time
}

func InsertTransaction(ctx context.Context, db *pgxpool.Pool, tx *Transaction) (string, error) {
	start := time.Now()

	var id string
	err := db.QueryRow(ctx,
		`INSERT INTO transactions (amount, category, merchant, transaction_date)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		tx.Amount, tx.Category, tx.Merchant, tx.TransactionDate).Scan(&id)
	if err != nil {
		return "", err
	}

	metrics.DbQueryDuration.WithLabelValues("insert_transaction").Observe(time.Since(start).Seconds())

	return id, nil
}
