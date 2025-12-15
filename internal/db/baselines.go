package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Baseline struct {
	Category     string
	AvgAmount    float64
	StdDevAmount float64
}

func UpdateBaselines(ctx context.Context, db *pgxpool.Pool, windowDays int) error {
	if windowDays <= 0 {
		return fmt.Errorf("windowDays must be > 0")
	}

	rows, err := db.Query(ctx,
		`SELECT category,
				COALESCE(AVG(amount), 0)::float8 AS average_amount,
				COALESCE(STDDEV(amount), 0)::float8 AS stddev_amount
		 FROM transactions
		 WHERE transaction_date >= NOW() - ($1::int || ' days')::interval
		 GROUP BY category`, windowDays)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var b Baseline
		if err := rows.Scan(&b.Category, &b.AvgAmount, &b.StdDevAmount); err != nil {
			return err
		}

		_, err := db.Exec(ctx,
			`INSERT INTO baselines (category, average_amount, stddev_amount, updated_at)
			 VALUES ($1, $2, $3, now())
			 ON CONFLICT (category) DO UPDATE
			 SET average_amount = EXCLUDED.average_amount,
			     stddev_amount = EXCLUDED.stddev_amount,
				 updated_at = now()`,
			b.Category, b.AvgAmount, b.StdDevAmount)
		if err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}
