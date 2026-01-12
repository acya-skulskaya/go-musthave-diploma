package balance

import (
	"errors"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/logger"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/middleware"
	balanceRepo "github.com/acya-skulskaya/go-musthave-diploma/internal/repository/balance"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/request"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/response"
	balanceService "github.com/acya-skulskaya/go-musthave-diploma/internal/service/balance"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/validator"
	"go.uber.org/zap"
	"net/http"
)

func PostAPIUserBalanceWithdraw(bs *balanceService.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceID, ok := r.Context().Value(middleware.ContextKeyTraceID).(string)
		if !ok {
			traceID = ""
		}

		withdrawal, err := request.Decode[request.Withdrawal](r)
		if err != nil {
			logger.Log.Debug("could not decode request",
				zap.String(middleware.TraceID, traceID),
				zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		problems, err := validator.IsValid(withdrawal)
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

		err = bs.Withdraw(r.Context(), withdrawal.Order, withdrawal.Sum)
		if err != nil {
			if errors.Is(err, balanceRepo.ErrNotEnoughBalanceToWithdraw) {
				http.Error(w, http.StatusText(http.StatusPaymentRequired), http.StatusPaymentRequired)
				return
			} else if errors.Is(err, balanceRepo.ErrWithdrawalForThisOrderAlreadyExists) {
				http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
				return
			}
			logger.Log.Error("could not withdraw",
				zap.String(middleware.TraceID, traceID),
				zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
