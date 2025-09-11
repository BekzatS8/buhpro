package repository

import (
	"context"
	"time"

	"github.com/BekzatS8/buhpro/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RefreshTokenRepo manages stored refresh tokens (hashes)
type RefreshTokenRepo interface {
	Create(ctx context.Context, t *models.RefreshToken) error
	GetByHash(ctx context.Context, hash string) (*models.RefreshToken, error)
	DeleteByUser(ctx context.Context, userID string) error
	DeleteByHash(ctx context.Context, hash string) error
}

type pgRefreshRepo struct {
	db *pgxpool.Pool
}

func NewRefreshRepo(db *pgxpool.Pool) RefreshTokenRepo {
	return &pgRefreshRepo{db: db}
}
func (r *pgRefreshRepo) Create(ctx context.Context, t *models.RefreshToken) error {
	_ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := r.db.Exec(_ctx, `
INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
VALUES ($1,$2,$3,$4,$5)
`, t.ID, t.UserID, t.TokenHash, t.ExpiresAt, t.CreatedAt)
	return err
}

func (r *pgRefreshRepo) GetByHash(ctx context.Context, hash string) (*models.RefreshToken, error) {
	_ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	t := &models.RefreshToken{}
	row := r.db.QueryRow(_ctx, `SELECT id,user_id,token_hash,expires_at,created_at FROM refresh_tokens WHERE token_hash=$1`, hash)
	if err := row.Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.CreatedAt); err != nil {
		return nil, err
	}
	return t, nil
}

func (r *pgRefreshRepo) DeleteByUser(ctx context.Context, userID string) error {
	_ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := r.db.Exec(_ctx, `DELETE FROM refresh_tokens WHERE user_id=$1`, userID)
	return err
}

func (r *pgRefreshRepo) DeleteByHash(ctx context.Context, hash string) error {
	_ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := r.db.Exec(_ctx, `DELETE FROM refresh_tokens WHERE token_hash=$1`, hash)
	return err
}
