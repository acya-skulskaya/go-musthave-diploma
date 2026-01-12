package balance

import (
	"github.com/acya-skulskaya/go-musthave-diploma/internal/logger"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/middleware"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/response"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/service/balance"
	"go.uber.org/zap"
	"net/http"
)

func GetAPIUserBalance(bs *balance.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceID, ok := r.Context().Value(middleware.ContextKeyTraceID).(string)
		if !ok {
			traceID = ""
		}

		userBalance, err := bs.Get(r.Context())
		if err != nil {
			logger.Log.Error("could not get user balance",
				zap.String(middleware.TraceID, traceID),
				zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		balanceResp := response.UserBalance{
			Current:   userBalance.Current,
			Withdrawn: userBalance.Withdrawn,
		}

		if err := response.Encode(w, http.StatusOK, balanceResp); err != nil {
			logger.Log.Error(logger.ErrorEncodingResponse,
				zap.String(middleware.TraceID, traceID),
				zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
