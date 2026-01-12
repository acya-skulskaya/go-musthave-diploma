package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/config/handlers"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	httpClient *http.Client
	Config     handlers.Config
	RetryAfter time.Duration
}

func New(config handlers.Config) *Service {
	return &Service{
		Config:     config,
		RetryAfter: 60 * time.Second,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *Service) GetOrderStatus(ctx context.Context, orderNumber string) (OrderResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/orders/%s", s.Config.AccrualAddress, orderNumber), http.NoBody)
	if err != nil {
		return OrderResponse{}, fmt.Errorf("could not create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return OrderResponse{}, fmt.Errorf("could not do request: %w", err)
	}
	//nolint:errcheck //возникает другая ошибка линтера
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return OrderResponse{}, fmt.Errorf("could not read response body: %w", err)
		}

		var res OrderResponse
		if err := json.Unmarshal(body, &res); err != nil {
			return OrderResponse{}, fmt.Errorf("could not decode response: %w", err)
		}
		return res, nil

	case http.StatusTooManyRequests:
		if hdr := resp.Header.Get("Retry-After"); hdr != "" {
			if sec, err := strconv.Atoi(strings.TrimSpace(hdr)); err == nil {
				s.RetryAfter = time.Duration(sec) * time.Second
			}
		}

		return OrderResponse{}, ErrTooManyRequests

	case http.StatusNoContent:
		return OrderResponse{}, ErrNoContent

	default:
		return OrderResponse{}, fmt.Errorf("no action for staus code %d", resp.StatusCode)
	}
}
