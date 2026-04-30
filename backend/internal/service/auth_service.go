package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"product-management/backend/internal/auth"
	"product-management/backend/internal/repository"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
	secret   string
	expiry   time.Duration
}

func NewAuthService(userRepo repository.UserRepository, secret string, expiry time.Duration) AuthService {
	return &authService{userRepo: userRepo, secret: secret, expiry: expiry}
}

func (s *authService) Login(ctx context.Context, username, password string) (string, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return "", fmt.Errorf("username dan password wajib diisi")
	}

	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := auth.GenerateToken(s.secret, user.ID, user.Username, s.expiry)
	if err != nil {
		return "", err
	}
	return token, nil
}
