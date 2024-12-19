package order

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostOrder(ctx context.Context, accountId string, products []OrderProduct) (*Order, error)
	GetOrderById(ctx context.Context, id string) (*Order, error)
	GetAccountOrders(ctx context.Context, accountId string) ([]Order, error)
}

type Order struct {
	ID            string         `json:"id"`
	Price         float64        `json:"price"`
	AccountID     string         `json:"accountId"`
	CreatedAt     time.Time      `json:"createdAt"`
	OrderProducts []OrderProduct `json:"products"`
}

type OrderProduct struct {
	ID          string  `json:"id"`
	Price       float64 `json:"price"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
}

type RequestedOrderProduct struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type OrderService struct {
	repository Repository
}

func NewOrderService(r Repository) *OrderService {
	return &OrderService{r}
}

func (s *OrderService) PostOrder(ctx context.Context, accountId string, products []OrderProduct) (*Order, error) {

	totoalPrice := 0.0
	for _, product := range products {
		totoalPrice += product.Price * float64(product.Quantity)
	}

	return s.repository.PostOrder(ctx, &Order{
		ID:            ksuid.New().String(),
		CreatedAt:     time.Now().UTC(),
		AccountID:     accountId,
		OrderProducts: products,
		Price:         totoalPrice,
	})
}
func (s *OrderService) GetOrderById(ctx context.Context, id string) (*Order, error) {
	return s.repository.GetOrderById(ctx, id)

}
func (s *OrderService) GetAccountOrders(ctx context.Context, accountId string) ([]Order, error) {
	return s.repository.GetAccountOrders(ctx, accountId)
}
