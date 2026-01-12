package auth

import (
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"net/http"
	"time"
)

func (s *Service) Login(w http.ResponseWriter, user models.User) error {
	token, err := BuildJWTString(s.SecretKey, user.ID)
	if err != nil {
		return fmt.Errorf("could not create token string: %w", err)
	}

	authCookie := http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(time.Hour * 24 * 365),
	}

	http.SetCookie(w, &authCookie)

	return nil
}
