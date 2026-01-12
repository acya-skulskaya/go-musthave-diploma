package auth

import (
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/repository/user"
	"golang.org/x/crypto/bcrypt"
)

const (
	ContextKeyUser = "user"
	CookieName     = "auth"
)

type Service struct {
	Repo      user.UserRepositoryInterface
	SecretKey string
}

func New(secretKey string, repo user.UserRepositoryInterface) *Service {
	return &Service{
		SecretKey: secretKey,
		Repo:      repo,
	}
}

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("could not generate hash: %w", err)
	}
	return string(hashedBytes), nil
}

func VerifyPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("could not compare hash and password: %w", err)
	}
	return nil
}
