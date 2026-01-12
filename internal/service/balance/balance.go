package balance

import "github.com/acya-skulskaya/go-musthave-diploma/internal/repository/balance"

type Service struct {
	Repo balance.BalanceRepositoryInterface
}

func New(repo balance.BalanceRepositoryInterface) *Service {
	return &Service{
		Repo: repo,
	}
}
