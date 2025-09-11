package repository

import (
	"context"
	"time"

	"github.com/BekzatS8/buhpro/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo interface {
	Create(u *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id string) (*models.User, error)
	Count() (int, error)
	Update(u *models.User) error
}

type pgUserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) UserRepo {
	return &pgUserRepo{db: db}
}

func (r *pgUserRepo) Create(u *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := r.db.Exec(ctx, `
		INSERT INTO users (id, email, phone, password_hash, full_name, role, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`, u.ID, u.Email, u.Phone, u.PasswordHash, u.FullName, u.Role, u.Status)
	return err
}

func (r *pgUserRepo) GetByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	u := &models.User{}
	row := r.db.QueryRow(ctx, `SELECT id,email,phone,full_name,role,status,password_hash,created_at,updated_at FROM users WHERE email=$1`, email)
	if err := row.Scan(&u.ID, &u.Email, &u.Phone, &u.FullName, &u.Role, &u.Status, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *pgUserRepo) GetByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	u := &models.User{}
	row := r.db.QueryRow(ctx, `SELECT id,email,phone,full_name,role,status,password_hash,created_at,updated_at FROM users WHERE id=$1`, id)
	if err := row.Scan(&u.ID, &u.Email, &u.Phone, &u.FullName, &u.Role, &u.Status, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *pgUserRepo) Count() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cnt int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&cnt); err != nil {
		return 0, err
	}
	return cnt, nil
}

func (r *pgUserRepo) Update(u *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, `
		UPDATE users
		SET full_name = $1,
			phone = $2,
			metadata = $3,
			updated_at = now()
		WHERE id = $4
	`, u.FullName, u.Phone, u.Metadata, u.ID)
	return err
}
