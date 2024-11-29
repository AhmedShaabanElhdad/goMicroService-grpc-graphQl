package catalog

import (
	"context"

	"github.com/segmentio/ksuid"
)

type catalogService struct {
	repository Repository
}

type Service interface {
	Close()
	PostCatalog(ctx context.Context, name, description string, price float64) error
	GetProducts(ctx context.Context, taken uint64, skip uint64) ([]Catalog, error)
	GetProductsByID(ctx context.Context, id string) (*Catalog, error)
	GetProductsWithIds(ctx context.Context, ids []string) ([]Catalog, error)
	SearchProduct(ctx context.Context, query string, taken uint64, skip uint64) ([]Catalog, error)
}

type Catalog struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func NewService(r Repository) Service {
	return &catalogService{r}
}

func (s catalogService) Close() {
	s.repository.Close()
}
func (s catalogService) PostCatalog(ctx context.Context, name, description string, price float64) error {
	catalog := Catalog{
		ID:          ksuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
	}
	return s.repository.PostCatalog(ctx, catalog)
}
func (s catalogService) GetProducts(ctx context.Context, taken uint64, skip uint64) ([]Catalog, error) {
	return s.repository.GetProducts(ctx, taken, skip)
}
func (s catalogService) GetProductsByID(ctx context.Context, id string) (*Catalog, error) {
	return s.repository.GetProductsByID(ctx, id)
}
func (s catalogService) GetProductsWithIds(ctx context.Context, ids []string) ([]Catalog, error) {
	return s.repository.GetProductsWithIds(ctx, ids)
}
func (s catalogService) SearchProduct(ctx context.Context, query string, taken uint64, skip uint64) ([]Catalog, error) {
	if taken > 100 || (skip == 0 && taken == 0) {
		taken = 100
	}
	return s.repository.SearchProduct(ctx, query, taken, skip)
}
