package repository

import (
	"context"
	"time"

	"github.com/BekzatS8/buhpro/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BidRepo interface {
	Create(ctx context.Context, b *domain.Bid) error
	GetByID(ctx context.Context, id string) (*domain.Bid, error)
	ListByOrder(ctx context.Context, orderID string) ([]*domain.Bid, error)
	Delete(ctx context.Context, id string) error
	MarkPaid(ctx context.Context, id string, paidAt time.Time) error
}

type pgBidRepo struct {
	db *pgxpool.Pool
}

func NewBidRepo(db *pgxpool.Pool) BidRepo { return &pgBidRepo{db: db} }

func (r *pgBidRepo) Create(ctx context.Context, b *domain.Bid) error {
	q := `INSERT INTO bids (id, order_id, executor_id, cover_text, price, proposed_deadline, attachments, status, metadata)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING created_at, updated_at`
	return r.db.QueryRow(ctx, q,
		b.ID, b.OrderID, b.ExecutorID, b.CoverText, b.Price, b.ProposedDeadline, b.Attachments, b.Status, b.Metadata,
	).Scan(&b.CreatedAt, &b.UpdatedAt)
}

func (r *pgBidRepo) GetByID(ctx context.Context, id string) (*domain.Bid, error) {
	b := &domain.Bid{}
	q := `SELECT id,order_id,executor_id,cover_text,price,proposed_deadline,attachments,status,paid_at,visibility_to_client,metadata,created_at,updated_at FROM bids WHERE id=$1`
	if err := r.db.QueryRow(ctx, q, id).Scan(
		&b.ID, &b.OrderID, &b.ExecutorID, &b.CoverText, &b.Price, &b.ProposedDeadline, &b.Attachments, &b.Status, &b.PaidAt, &b.VisibleToClient, &b.Metadata, &b.CreatedAt, &b.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return b, nil
}

func (r *pgBidRepo) ListByOrder(ctx context.Context, orderID string) ([]*domain.Bid, error) {
	rows, err := r.db.Query(ctx, `SELECT id,order_id,executor_id,cover_text,price,proposed_deadline,attachments,status,paid_at,visibility_to_client,metadata,created_at,updated_at FROM bids WHERE order_id=$1`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Bid
	for rows.Next() {
		b := &domain.Bid{}
		if err := rows.Scan(&b.ID, &b.OrderID, &b.ExecutorID, &b.CoverText, &b.Price, &b.ProposedDeadline, &b.Attachments, &b.Status, &b.PaidAt, &b.VisibleToClient, &b.Metadata, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, nil
}

func (r *pgBidRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM bids WHERE id=$1 AND status IN ('created','pending_payment')`, id)
	return err
}

func (r *pgBidRepo) MarkPaid(ctx context.Context, id string, paidAt time.Time) error {
	_, err := r.db.Exec(ctx, `UPDATE bids SET status='paid', paid_at=$1, visibility_to_client=true, updated_at=now() WHERE id=$2`, paidAt, id)
	return err
}
