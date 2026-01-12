package order

import (
	"github.com/acya-skulskaya/go-musthave-diploma/internal/repository/order"
)

type Service struct {
	Repo order.OrderRepositoryInterface
}

func New(repo order.OrderRepositoryInterface) *Service {
	return &Service{
		Repo: repo,
	}
}
