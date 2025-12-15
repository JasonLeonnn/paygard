package services

import (
	"context"

	"github.com/JasonLeonnn/paygard/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlertService struct {
	pool *pgxpool.Pool
}

func NewAlertService(pool *pgxpool.Pool) *AlertService {
	return &AlertService{pool: pool}
}

func (s *AlertService) GetAlerts(ctx context.Context, limit int) ([]db.Alert, error) {
	return db.ListAlerts(ctx, s.pool, limit)
}
