package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl/catalog"
	"github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl/order"
)

var ErrInvalidParameter = errors.New("invalid parameter")

type mutationResolver struct {
	server *Server
}

func (resolver *mutationResolver) CreateAccount(ctx context.Context, accountInput AccountInput) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	account, err := resolver.server.accountClient.PostAccount(ctx, accountInput.Name)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Account{
		ID:   account.ID,
		Name: account.Name,
	}, nil
}

func (resolver *mutationResolver) CreateProduct(ctx context.Context, productInput ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	product, err := resolver.server.catalogClient.PostCatalog(ctx, catalog.Catalog{
		Name:        productInput.Name,
		Description: productInput.Description,
		Price:       productInput.Price,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}, nil
}

func (resolver *mutationResolver) CreateOrder(ctx context.Context, orderInput OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var orderProducts []order.RequestedOrderProduct
	for _, productInput := range orderInput.Products {
		if productInput.Quantity <= 0 {
			return nil, ErrInvalidParameter
		}

		orderProducts = append(orderProducts, order.RequestedOrderProduct{
			ID:       productInput.ID,
			Quantity: productInput.Quantity,
		})

	}
	createdOrder, err := resolver.server.orderClient.CreateOrder(ctx, orderProducts, orderInput.AccountID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var products []*OrderProduct
	for _, product := range createdOrder.OrderProducts {
		products = append(products, &OrderProduct{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Quantity:    product.Quantity,
		})
	}

	return &Order{
		ID:         createdOrder.ID,
		CreatedAt:  createdOrder.CreatedAt.String(),
		TotlaPrice: createdOrder.Price,
		Products:   products,
	}, nil
}

func (resolver *mutationResolver) Login(ctx context.Context, email string, password string) (*User, error) {
	return nil, nil
}

func (resolver *mutationResolver) Register(ctx context.Context, input NewUser) (*User, error) {
	return nil, nil
}
