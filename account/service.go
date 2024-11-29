package account

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostAccoun(ctx context.Context, name string) (*Account, error)
	GetAccount(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type Account struct {
	ID   string `json:"Id"`
	Name string `json:"name"`
}

type accounService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &accounService{
		repository: r,
	}
}

func (s *accounService) PostAccoun(ctx context.Context, name string) (*Account, error) {
	a := &Account{
		Name: name,
		ID:   ksuid.New().String(),
	}
	if err := s.repository.AddAccount(ctx, *a); err != nil {
		return nil, err
	}
	return a, nil

}
func (s *accounService) GetAccount(ctx context.Context, id string) (*Account, error) {
	return s.repository.GetAccountByID(ctx, id)
}
func (s *accounService) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.FetchAccounts(ctx, skip, take)
}
