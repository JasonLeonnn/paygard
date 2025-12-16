package services

import (
	"context"
	"log"

	"github.com/JasonLeonnn/jalytics/internal/db"
	"github.com/JasonLeonnn/jalytics/internal/metrics"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionService struct {
	pool *pgxpool.Pool
}

func NewTransactionService(pool *pgxpool.Pool) *TransactionService {
	return &TransactionService{pool: pool}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, tx *db.Transaction) (bool, string, error) {
	id, err := db.InsertTransaction(ctx, s.pool, tx)
	if err != nil {
		return false, "", err
	}
	metrics.TransactionCounter.WithLabelValues(tx.Category, tx.Merchant).Inc()
	tx.ID = id

	isAnomaly, severity, err := db.CheckAnomaly(ctx, s.pool, tx, false)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("No baseline found for category=%s, skipping anomaly detection. Run baseline update to create baselines.", tx.Category)
			isAnomaly = false
		} else {
			log.Printf("Error checking anomaly: %v", err)
			return false, "", err
		}
	} else if isAnomaly {
		metrics.AnomalyCounter.Inc()
		if severity != "" {
			metrics.AnomalyBySeverityCounter.WithLabelValues(severity).Inc()
		}
		log.Printf("anomaly detected: tx_id=%s category=%s severity=%s amount=%.2f", tx.ID, tx.Category, severity, tx.Amount)
	}

	return isAnomaly, id, nil
}
