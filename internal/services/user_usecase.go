package services

import (
	"context"
	"errors"
	"time"

	"github.com/BekzatS8/buhpro/internal/models"
	"github.com/BekzatS8/buhpro/internal/repository"
	"github.com/BekzatS8/buhpro/pkg/auth"
	"github.com/google/uuid"
)

type UserUsecase struct {
	repo        repository.UserRepo
	refreshRepo repository.RefreshTokenRepo
	jwtSecret   string
	jwtTTL      int
	refreshTTL  int // days
}

func NewUserUsecase(r repository.UserRepo, rr repository.RefreshTokenRepo, jwtSecret string, jwtTTL int, refreshTTLDays int) *UserUsecase {
	return &UserUsecase{repo: r, refreshRepo: rr, jwtSecret: jwtSecret, jwtTTL: jwtTTL, refreshTTL: refreshTTLDays}
}

// UserUpdate — DTO для обновления профиля
type UserUpdate struct {
	FullName *string                `json:"full_name,omitempty"`
	Phone    *string                `json:"phone,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Register now returns access and refresh tokens
func (uc *UserUsecase) Register(email, phone, fullName, password string, role string) (string, string, error) {
	if u, _ := uc.repo.GetByEmail(email); u != nil {
		return "", "", errors.New("email already registered")
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return "", "", err
	}
	user := &models.User{
		ID:           uuid.NewString(),
		Email:        email,
		Phone:        phone,
		FullName:     fullName,
		Role:         role,
		Status:       "active",
		PasswordHash: hash,
	}
	if err := uc.repo.Create(user); err != nil {
		return "", "", err
	}
	// generate tokens
	access, refresh, err := auth.GenerateTokens(uc.jwtSecret, user.ID, user.Role, uc.jwtTTL, uc.refreshTTL)
	if err != nil {
		return "", "", err
	}
	// store refresh token (hashed)
	rt := &models.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		TokenHash: auth.HashToken(refresh),
		ExpiresAt: time.Now().Add(time.Hour * 24 * time.Duration(uc.refreshTTL)),
		CreatedAt: time.Now(),
	}
	if err := uc.refreshRepo.Create(context.Background(), rt); err != nil {
		// return error to be safe
		return "", "", err
	}
	return access, refresh, nil
}

// Login returns access and refresh tokens and saves refresh
func (uc *UserUsecase) Login(email, password string) (string, string, error) {
	u, err := uc.repo.GetByEmail(email)
	if err != nil {
		return "", "", err
	}
	if err := auth.CheckPassword(u.PasswordHash, password); err != nil {
		return "", "", errors.New("invalid credentials")
	}
	access, refresh, err := auth.GenerateTokens(uc.jwtSecret, u.ID, u.Role, uc.jwtTTL, uc.refreshTTL)
	if err != nil {
		return "", "", err
	}
	// delete previous refresh tokens for this user (optional)
	_ = uc.refreshRepo.DeleteByUser(context.Background(), u.ID)
	rt := &models.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    u.ID,
		TokenHash: auth.HashToken(refresh),
		ExpiresAt: time.Now().Add(time.Hour * 24 * time.Duration(uc.refreshTTL)),
		CreatedAt: time.Now(),
	}
	if err := uc.refreshRepo.Create(context.Background(), rt); err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func (uc *UserUsecase) RepoCount() (int, error) {
	return uc.repo.Count()
}

func (uc *UserUsecase) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := auth.ParseRefreshToken(uc.jwtSecret, refreshToken)
	if err != nil {
		return "", "", err
	}
	// check stored hash
	hash := auth.HashToken(refreshToken)
	stored, err := uc.refreshRepo.GetByHash(ctx, hash)
	if err != nil {
		return "", "", err
	}
	if stored.ExpiresAt.Before(time.Now()) {
		return "", "", errors.New("refresh token expired")
	}
	// generate new tokens
	access, newRefresh, err := auth.GenerateTokens(uc.jwtSecret, claims.UserID, claims.Role, uc.jwtTTL, uc.refreshTTL)
	if err != nil {
		return "", "", err
	}
	// delete old and save new
	if err := uc.refreshRepo.DeleteByHash(ctx, hash); err != nil {
		return "", "", err
	}
	rt := &models.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    claims.UserID,
		TokenHash: auth.HashToken(newRefresh),
		ExpiresAt: time.Now().Add(time.Hour * 24 * time.Duration(uc.refreshTTL)),
		CreatedAt: time.Now(),
	}
	if err := uc.refreshRepo.Create(ctx, rt); err != nil {
		return "", "", err
	}
	return access, newRefresh, nil
}

// Logout removes all refresh tokens for a user (or a single token if needed)
func (uc *UserUsecase) Logout(ctx context.Context, userID string) error {
	return uc.refreshRepo.DeleteByUser(ctx, userID)
}

// GetProfile and UpdateProfile
func (uc *UserUsecase) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	u, err := uc.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	u.PasswordHash = ""
	return u, nil
}

func (uc *UserUsecase) UpdateProfile(ctx context.Context, userID string, upd UserUpdate) (*models.User, error) {
	u, err := uc.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if upd.FullName != nil {
		u.FullName = *upd.FullName
	}
	if upd.Phone != nil {
		u.Phone = *upd.Phone
	}
	if upd.Metadata != nil {
		u.Metadata = upd.Metadata
	}
	if err := uc.repo.Update(u); err != nil {
		return nil, err
	}
	u.PasswordHash = ""
	return u, nil
}
