package response

import (
	"time"
)

type Withdrawal struct {
	ProcessedAt *time.Time `json:"processed_at"`
	OrderNumber string     `json:"order"`
	Sum         float64    `json:"sum"`
}
