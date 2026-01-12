package middleware

import (
	"context"
	"errors"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/config"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/logger"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	repoUser "github.com/acya-skulskaya/go-musthave-diploma/internal/repository/user"
	serviceAuth "github.com/acya-skulskaya/go-musthave-diploma/internal/service/auth"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

const (
	ContextKeyUser ContextKey = "user"
)

func HandleCookieAuth(cfg *config.Config, authService *serviceAuth.Service) func(next http.Handler) http.Handler {
	cookieAuth := NewCookieAuth(cfg, authService)
	return cookieAuth.Handler
}

type CookieAUth struct {
	AuthService *serviceAuth.Service
	SecretKey   string
}

func NewCookieAuth(cfg *config.Config, authService *serviceAuth.Service) *CookieAUth {
	c := &CookieAUth{SecretKey: cfg.Auth.SecretKey, AuthService: authService}
	return c
}

func (c *CookieAUth) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID, ok := r.Context().Value(ContextKeyTraceID).(string)
		if !ok {
			traceID = ""
		}

		cookieValue := ""
		cookie, err := r.Cookie(serviceAuth.CookieName)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				logger.Log.Debug("auth cookie not found",
					zap.String(TraceID, traceID))
			default:
				logger.Log.Debug("could not get cookie",
					zap.String(TraceID, traceID),
					zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			cookieValue = cookie.Value
		}

		if cookieValue == "" {
			logger.Log.Debug("auth cookie is empty",
				zap.String(TraceID, traceID))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, err := serviceAuth.GetUserIDFromToken(c.SecretKey, cookieValue)
		if err != nil {
			if errors.Is(err, serviceAuth.ErrTokenIsNotValid) {
				logger.Log.Debug("token is not valid",
					zap.String(TraceID, traceID))
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				logger.Log.Debug("error getting user id from auth token",
					zap.String(TraceID, traceID),
					zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		if userID == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		logger.Log.Debug("got user id",
			zap.String(TraceID, traceID),
			zap.Uint("userID", userID))

		user, err := UsersCont.get(userID)
		if err != nil {
			user, err = c.AuthService.Repo.Get(r.Context(), userID)
			if err != nil {
				logger.Log.Debug("error getting user",
					zap.String(TraceID, traceID),
					zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			UsersCont.add(userID, user)
		}

		ctx := context.WithValue(r.Context(), ContextKeyUser, user)
		r = r.WithContext(ctx)

		// передаём управление хендлеру
		next.ServeHTTP(w, r)
	})
}

// UsersContainer *****************************************************************************************
type UsersContainer struct {
	users map[uint]models.User
	mu    sync.RWMutex
}

func (c *UsersContainer) add(id uint, user models.User) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.users[id]; ok {
		return
	}

	c.users[id] = user
}

func (c *UsersContainer) get(id uint) (user models.User, err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.users[id]
	if !ok {
		return models.User{}, repoUser.ErrNotFound
	}

	return item, nil
}

var UsersCont = UsersContainer{users: make(map[uint]models.User)}
