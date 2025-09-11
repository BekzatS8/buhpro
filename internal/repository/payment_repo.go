package repository

import (
	"context"

	"github.com/BekzatS8/buhpro/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepo interface {
	Create(ctx context.Context, p *models.Payment) error
	GetByID(ctx context.Context, id string) (*models.Payment, error)
	UpdateStatus(ctx context.Context, id, status string) error
}

type pgPaymentRepo struct {
	db *pgxpool.Pool
}

func NewPaymentRepo(db *pgxpool.Pool) PaymentRepo { return &pgPaymentRepo{db: db} }

func (r *pgPaymentRepo) Create(ctx context.Context, p *models.Payment) error {
	// make placeholders count match columns (14)
	q := `INSERT INTO payments (
        id, user_id, organization_id, related_type, related_id,
        provider, provider_payment_id, amount, currency, status,
        items, idempotency_key, expires_at, webhook_meta
    ) VALUES (
        $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14
    ) RETURNING created_at, updated_at`

	// Convert possible empty string pointers to nil so PG gets NULL instead of "" (invalid uuid)
	var userID interface{}
	if p.UserID == nil || *p.UserID == "" {
		userID = nil
	} else {
		userID = *p.UserID
	}

	var orgID interface{}
	if p.OrganizationID == nil || *p.OrganizationID == "" {
		orgID = nil
	} else {
		orgID = *p.OrganizationID
	}

	var relatedID interface{}
	if p.RelatedID == nil || *p.RelatedID == "" {
		relatedID = nil
	} else {
		relatedID = *p.RelatedID
	}

	// items and webhook_meta are JSONB; pgx will accept []byte, map[string]interface{}, or nil
	var items interface{} = p.Items
	if items == nil {
		items = nil
	}

	var webhook interface{} = p.WebhookMeta
	if webhook == nil {
		webhook = nil
	}

	// execute
	return r.db.QueryRow(ctx, q,
		p.ID, userID, orgID, p.RelatedType, relatedID,
		p.Provider, p.ProviderPaymentID, p.Amount, p.Currency, p.Status,
		items, p.IdempotencyKey, p.ExpiresAt, webhook,
	).Scan(&p.CreatedAt, &p.UpdatedAt)
}

func (r *pgPaymentRepo) GetByID(ctx context.Context, id string) (*models.Payment, error) {
	p := &models.Payment{}
	q := `SELECT id,user_id,organization_id,related_type,related_id,provider,provider_payment_id,amount,currency,status,items,idempotency_key,expires_at,webhook_meta,created_at,updated_at FROM payments WHERE id=$1`
	if err := r.db.QueryRow(ctx, q, id).Scan(
		&p.ID, &p.UserID, &p.OrganizationID, &p.RelatedType, &p.RelatedID, &p.Provider, &p.ProviderPaymentID, &p.Amount, &p.Currency, &p.Status, &p.Items, &p.IdempotencyKey, &p.ExpiresAt, &p.WebhookMeta, &p.CreatedAt, &p.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return p, nil
}

func (r *pgPaymentRepo) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.Exec(ctx, `UPDATE payments SET status=$1, updated_at=now() WHERE id=$2`, status, id)
	return err
}
