package withdrawals

import (
	"errors"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/logger"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/middleware"
	balanceRepo "github.com/acya-skulskaya/go-musthave-diploma/internal/repository/balance"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/response"
	balanceService "github.com/acya-skulskaya/go-musthave-diploma/internal/service/balance"
	"go.uber.org/zap"
	"net/http"
)

func GetAPIUserWithdrawals(bs *balanceService.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceID, ok := r.Context().Value(middleware.ContextKeyTraceID).(string)
		if !ok {
			traceID = ""
		}

		userBalanceFlow, err := bs.WithdrawalsByUser(r.Context())
		if err != nil {
			if errors.Is(err, balanceRepo.ErrNoWithdrawals) {
				w.WriteHeader(http.StatusNoContent)
				return
			} else {
				logger.Log.Error("could not get user withdrawals",
					zap.String(middleware.TraceID, traceID),
					zap.Error(err))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		var withdrawals []response.Withdrawal

		for _, flow := range userBalanceFlow {
			withdrawal := response.Withdrawal{
				OrderNumber: flow.OrderNumber,
				Sum:         flow.Amount * -1,
				ProcessedAt: flow.ProcessedAt,
			}
			withdrawals = append(withdrawals, withdrawal)
		}

		if len(withdrawals) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if err := response.Encode(w, http.StatusOK, withdrawals); err != nil {
			logger.Log.Error(logger.ErrorEncodingResponse,
				zap.String(middleware.TraceID, traceID),
				zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
