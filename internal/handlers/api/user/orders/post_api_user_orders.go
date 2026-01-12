package orders

import (
	"errors"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/logger"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/middleware"
	orderRepo "github.com/acya-skulskaya/go-musthave-diploma/internal/repository/order"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/request"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/response"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/service/order"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/validator"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func PostAPIUserOrders(os *order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceID, ok := r.Context().Value(middleware.ContextKeyTraceID).(string)
		if !ok {
			traceID = ""
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var orderID = request.OrderID(string(body))
		problems, err := validator.IsValid(orderID)
		if err != nil {
			if errEnc := response.Encode(w, http.StatusUnprocessableEntity, problems); errEnc != nil {
				logger.Log.Error(logger.ErrorEncodingResponse,
					zap.String(middleware.TraceID, traceID),
					zap.Error(err))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			return
		}

		_, err = os.Store(r.Context(), string(orderID))
		if err != nil {
			//nolint:gocritic // неудобно
			if errors.Is(err, orderRepo.ErrOrderAlreadyExistsByCurrentUser) {
				http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
				return
			} else if errors.Is(err, orderRepo.ErrOrderAlreadyExistsFromAnotherUser) {
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			} else {
				logger.Log.Error("could not store order",
					zap.String(middleware.TraceID, traceID),
					zap.Error(err))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
