package orders

import (
	"github.com/acya-skulskaya/go-musthave-diploma/internal/logger"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/middleware"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/response"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/service/order"
	"go.uber.org/zap"
	"net/http"
)

func GetAPIUserOrders(os *order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceID, ok := r.Context().Value(middleware.ContextKeyTraceID).(string)
		if !ok {
			traceID = ""
		}

		orders, err := os.All(r.Context())
		if err != nil {
			logger.Log.Error("could not get all orders",
				zap.String(middleware.TraceID, traceID),
				zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if len(orders) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var ordersResponse []response.Order

		for _, order := range orders {
			orderResponse := response.Order{
				UploadedAt:  order.UploadedAt,
				OrderNumber: order.OrderNumber,
				Status:      order.Status,
				Accrual:     order.Accrual,
			}
			ordersResponse = append(ordersResponse, orderResponse)
		}

		if err := response.Encode(w, http.StatusOK, ordersResponse); err != nil {
			logger.Log.Error(logger.ErrorEncodingResponse,
				zap.String(middleware.TraceID, traceID),
				zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
