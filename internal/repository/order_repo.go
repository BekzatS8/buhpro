package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/BekzatS8/buhpro/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepo interface {
	Create(ctx context.Context, o *domain.Order) error
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	List(ctx context.Context, filters map[string]string, page, perPage int) ([]*domain.Order, int, error)
	Update(ctx context.Context, o *domain.Order) error
	Delete(ctx context.Context, id string) error
	SetStatus(ctx context.Context, id, status string) error
	SelectExecutor(ctx context.Context, orderID, bidID string) error
	AddHistory(ctx context.Context, actorID, action, objectType, objectID string, payload map[string]interface{}) error
}

type pgOrderRepo struct {
	db *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) OrderRepo {
	return &pgOrderRepo{db: db}
}

func (r *pgOrderRepo) Create(ctx context.Context, o *domain.Order) error {
	query := `INSERT INTO orders (id, org_id, client_user_id, title, description, category, subcategory, region,
		mode_online, deadline, budget_min, budget_max, currency, status, promotion_flags, attachments, chosen_bid_id)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
	RETURNING created_at, updated_at`
	err := r.db.QueryRow(ctx, query,
		o.ID, o.OrgID, o.ClientUserID, o.Title, o.Description, o.Category, o.Subcategory, o.Region,
		o.ModeOnline, o.Deadline, o.BudgetMin, o.BudgetMax, o.Currency, o.Status, o.Promotion, o.Attachments, o.ChosenBidID,
	).Scan(&o.CreatedAt, &o.UpdatedAt)
	return err
}

func (r *pgOrderRepo) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	o := &domain.Order{}
	query := `SELECT id, org_id, client_user_id, title, description, category, subcategory, region, mode_online,
		deadline, budget_min, budget_max, currency, status, promotion_flags, attachments, chosen_bid_id, created_at, published_at, updated_at
		FROM orders WHERE id=$1`
	if err := r.db.QueryRow(ctx, query, id).Scan(
		&o.ID, &o.OrgID, &o.ClientUserID, &o.Title, &o.Description, &o.Category, &o.Subcategory, &o.Region, &o.ModeOnline,
		&o.Deadline, &o.BudgetMin, &o.BudgetMax, &o.Currency, &o.Status, &o.Promotion, &o.Attachments, &o.ChosenBidID,
		&o.CreatedAt, &o.PublishedAt, &o.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return o, nil
}

func (r *pgOrderRepo) List(ctx context.Context, filters map[string]string, page, perPage int) ([]*domain.Order, int, error) {
	var where []string
	var args []interface{}
	i := 1

	// filters: status, category, region, min_budget, max_budget, top, pinned
	if v, ok := filters["status"]; ok && v != "" {
		where = append(where, fmt.Sprintf("status = $%d", i))
		args = append(args, v)
		i++
	}
	if v, ok := filters["category"]; ok && v != "" {
		where = append(where, fmt.Sprintf("category = $%d", i))
		args = append(args, v)
		i++
	}
	if v, ok := filters["region"]; ok && v != "" {
		where = append(where, fmt.Sprintf("region = $%d", i))
		args = append(args, v)
		i++
	}
	if v, ok := filters["min_budget"]; ok && v != "" {
		where = append(where, fmt.Sprintf("budget_min >= $%d", i))
		args = append(args, v)
		i++
	}
	if v, ok := filters["max_budget"]; ok && v != "" {
		where = append(where, fmt.Sprintf("budget_max <= $%d", i))
		args = append(args, v)
		i++
	}
	// TODO: top/pinned use promotion_flags JSONB fields if needed

	q := "SELECT id, org_id, client_user_id, title, description, category, subcategory, region, mode_online, deadline, budget_min, budget_max, currency, status, promotion_flags, attachments, chosen_bid_id, created_at, published_at, updated_at FROM orders"
	if len(where) > 0 {
		q += " WHERE " + strings.Join(where, " AND ")
	}
	// count
	countQ := "SELECT count(*) FROM orders"
	if len(where) > 0 {
		countQ += " WHERE " + strings.Join(where, " AND ")
	}
	var total int
	if err := r.db.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// pagination & ordering
	offset := (page - 1) * perPage
	q += " ORDER BY created_at DESC LIMIT $%d OFFSET $%d"
	args = append(args, perPage, offset)
	// replace placeholders
	q = fmt.Sprintf(q, len(args)-1, len(args))

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*domain.Order
	for rows.Next() {
		o := &domain.Order{}
		if err := rows.Scan(
			&o.ID, &o.OrgID, &o.ClientUserID, &o.Title, &o.Description, &o.Category, &o.Subcategory, &o.Region, &o.ModeOnline,
			&o.Deadline, &o.BudgetMin, &o.BudgetMax, &o.Currency, &o.Status, &o.Promotion, &o.Attachments, &o.ChosenBidID,
			&o.CreatedAt, &o.PublishedAt, &o.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		out = append(out, o)
	}
	return out, total, nil
}

func (r *pgOrderRepo) Update(ctx context.Context, o *domain.Order) error {
	query := `UPDATE orders SET title=$1, description=$2, category=$3, subcategory=$4, region=$5, mode_online=$6,
		deadline=$7, budget_min=$8, budget_max=$9, currency=$10, promotion_flags=$11, attachments=$12, updated_at=now()
		WHERE id=$13 RETURNING updated_at`
	return r.db.QueryRow(ctx, query,
		o.Title, o.Description, o.Category, o.Subcategory, o.Region, o.ModeOnline,
		o.Deadline, o.BudgetMin, o.BudgetMax, o.Currency, o.Promotion, o.Attachments, o.ID,
	).Scan(&o.UpdatedAt)
}

func (r *pgOrderRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM orders WHERE id=$1`, id)
	return err
}

func (r *pgOrderRepo) SetStatus(ctx context.Context, id, status string) error {
	_, err := r.db.Exec(ctx, `UPDATE orders SET status=$1, updated_at=now() WHERE id=$2`, status, id)
	return err
}

func (r *pgOrderRepo) SelectExecutor(ctx context.Context, orderID, bidID string) error {
	_, err := r.db.Exec(ctx, `UPDATE orders SET chosen_bid_id=$1, status='executor_selected', updated_at=now() WHERE id=$2`, bidID, orderID)
	return err
}

func (r *pgOrderRepo) AddHistory(ctx context.Context, actorID, action, objectType, objectID string, payload map[string]interface{}) error {
	q := `INSERT INTO audit_logs (actor_id, action, object_type, object_id, payload) VALUES ($1,$2,$3,$4,$5)`
	_, err := r.db.Exec(ctx, q, actorID, action, objectType, objectID, payload)
	return err
}
