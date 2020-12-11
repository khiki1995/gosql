package security

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)
type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (s *Service) Auth(ctx context.Context, login, password string) (ok bool) {
	var managerID int64
	err := s.pool.QueryRow(ctx, `SELECT id FROM managers WHERE login = $1 and password = $2`, login, password).Scan(&managerID)
	if err != nil {
		return false
	}
	return true
}