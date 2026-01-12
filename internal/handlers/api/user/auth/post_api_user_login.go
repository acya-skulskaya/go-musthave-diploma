package auth

import (
	"errors"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/logger"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/middleware"
	repoUser "github.com/acya-skulskaya/go-musthave-diploma/internal/repository/user"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/request"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/response"
	serviceAuth "github.com/acya-skulskaya/go-musthave-diploma/internal/service/auth"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/validator"
	"go.uber.org/zap"
	"net/http"
)

func PostAPIUserLogin(as *serviceAuth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceID, ok := r.Context().Value(middleware.ContextKeyTraceID).(string)
		if !ok {
			traceID = ""
		}

		userCreds, err := request.Decode[request.UserCredits](r)
		if err != nil {
			logger.Log.Debug("could not decode request",
				zap.String(middleware.TraceID, traceID),
				zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		problems, err := validator.IsValid(userCreds)
		if err != nil {
			if errEnc := response.Encode(w, http.StatusBadRequest, problems); errEnc != nil {
				logger.Log.Error(logger.ErrorEncodingResponse,
					zap.String(middleware.TraceID, traceID),
					zap.Error(err))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			return
		}

		user, err := as.Repo.Get(r.Context(), userCreds.Login)
		if err != nil {
			if errors.Is(err, repoUser.ErrNotFound) {
				logger.Log.Debug(repoUser.ErrMsgLoginNotFound,
					zap.String(middleware.TraceID, traceID),
					zap.Error(err))
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			} else {
				logger.Log.Error(repoUser.ErrMsgLoginNotFound,
					zap.String(middleware.TraceID, traceID),
					zap.Error(err))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		if err = serviceAuth.VerifyPassword(user.Password, userCreds.Password); err != nil {
			logger.Log.Debug("could not verify password",
				zap.String(middleware.TraceID, traceID),
				zap.Error(err))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		err = as.Login(w, user)
		if err != nil {
			logger.Log.Error(repoUser.ErrMsgCouldNotLogIn,
				zap.String(middleware.TraceID, traceID),
				zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
