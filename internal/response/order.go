package response

import (
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"time"
)

type Order struct {
	UploadedAt  *time.Time         `json:"uploaded_at"`
	OrderNumber string             `json:"number"`
	Status      models.OrderStatus `json:"status"`
	Accrual     float64            `json:"accrual,omitempty"`
}
