package repository

import (
	"context"
	"fmt"
	"kibit/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LeadRepo struct{ pool *pgxpool.Pool }

func NewLeadRepo(p *pgxpool.Pool) *LeadRepo { return &LeadRepo{p} }
func (r *LeadRepo) Create(ctx context.Context, req *model.CreateLeadRequest) (int64, error) {
	var id int64
	err := r.pool.QueryRow(ctx, `
        INSERT INTO leads (name, email, phone, company, app_type, message, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, NOW()) RETURNING id`,
		req.Name, req.Email, req.Phone, req.Company, req.AppType, req.Message).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert lead: %w", err)
	}
	return id, nil
}
