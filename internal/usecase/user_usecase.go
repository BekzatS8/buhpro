package usecase

import (
	"errors"
	"github.com/google/uuid"

	"github.com/BekzatS8/buhpro/internal/domain"
	"github.com/BekzatS8/buhpro/internal/repository"
	"github.com/BekzatS8/buhpro/pkg/auth"
)

type UserUsecase struct {
	repo      repository.UserRepo
	jwtSecret string
	jwtTTL    int
}

func NewUserUsecase(r repository.UserRepo, jwtSecret string, jwtTTL int) *UserUsecase {
	return &UserUsecase{repo: r, jwtSecret: jwtSecret, jwtTTL: jwtTTL}
}

func (uc *UserUsecase) Register(email, phone, fullName, password string, role string) (string, error) {
	// check exists
	if u, _ := uc.repo.GetByEmail(email); u != nil {
		return "", errors.New("email already registered")
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return "", err
	}
	user := &domain.User{
		ID:           uuid.NewString(),
		Email:        email,
		Phone:        phone,
		FullName:     fullName,
		Role:         role,
		Status:       "active",
		PasswordHash: hash,
	}
	if err := uc.repo.Create(user); err != nil {
		return "", err
	}
	// generate token
	token, err := auth.GenerateToken(uc.jwtSecret, user.ID, user.Role, uc.jwtTTL)
	return token, err
}

func (uc *UserUsecase) Login(email, password string) (string, error) {
	u, err := uc.repo.GetByEmail(email)
	if err != nil {
		return "", err
	}
	if err := auth.CheckPassword(u.PasswordHash, password); err != nil {
		return "", errors.New("invalid credentials")
	}
	return auth.GenerateToken(uc.jwtSecret, u.ID, u.Role, uc.jwtTTL)
}
func (uc *UserUsecase) RepoCount() (int, error) {
	return uc.repo.Count()
}
