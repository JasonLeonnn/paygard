package db

import (
	"context"
	"math"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func severityForZScore(z float64, strict bool) (string, string) {
	if strict {
		switch {
		case z >= 5:
			return "critical", "Amount exceeds 5\u03c3 baseline"
		case z >= 3:
			return "high", "Amount exceeds 3\u03c3 baseline"
		case z >= 2:
			return "medium", "Amount exceeds 2\u03c3 baseline"
		default:
			return "", ""
		}
	}

	switch {
	case z >= 5:
		return "critical", "Amount exceeds 5\u03c3 baseline"
	case z >= 3:
		return "high", "Amount exceeds 3\u03c3 baseline"
	default:
		return "", ""
	}
}

func CheckAnomaly(ctx context.Context, db *pgxpool.Pool, tx *Transaction, strict bool) (bool, string, error) {
	var avg, stddev float64
	err := db.QueryRow(ctx,
		`SELECT average_amount::float8, stddev_amount::float8
		 FROM baselines
		 WHERE category = $1`, tx.Category).Scan(&avg, &stddev)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, "", nil
		}
		return false, "", err
	}

	if stddev <= 0 {
		if tx.Amount >= avg*2 && tx.Amount-avg > 50 {
			_, err := db.Exec(ctx,
				`INSERT INTO alerts (transaction_id, alert_message, severity)
				 VALUES ($1, $2, $3)`,
				tx.ID, "Amount is more than 2x stable baseline", "high")
			if err != nil {
				return false, "", err
			}
			return true, "high", nil
		}
		return false, "", nil
	}

	z := (tx.Amount - avg) / stddev
	z = math.Round(z*100) / 100

	severity, message := severityForZScore(z, strict)
	if severity == "" {
		return false, "", nil
	}

	_, err = db.Exec(ctx,
		`INSERT INTO alerts (transaction_id, alert_message, severity)
			 VALUES ($1, $2, $3)`,
		tx.ID, message, severity)
	if err != nil {
		return false, "", err
	}
	return true, severity, nil
}
