package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Alert struct {
	ID            string
	TransactionID string
	AlertMessage  string
	Severity      string
	CreatedAt     time.Time
}

func ListAlerts(ctx context.Context, db *pgxpool.Pool, limit int) ([]Alert, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	rows, err := db.Query(ctx,
		`SELECT id, transaction_id, alert_message, severity, created_at
		 FROM alerts
		 ORDER BY created_at DESC
		 LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []Alert
	for rows.Next() {
		var a Alert
		if err := rows.Scan(&a.ID, &a.TransactionID, &a.AlertMessage, &a.Severity, &a.CreatedAt); err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return alerts, nil
}
