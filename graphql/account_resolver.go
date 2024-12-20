package main

import (
	"context"
	"log"
	"time"
)

type accountResolver struct {
	server *Server
}

func (r *accountResolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(3*time.Second))
	defer cancel()

	orderList, err := r.server.orderClient.GetAccountOrder(ctx, obj.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var orders []*Order
	for _, order := range orderList {

		var products []*OrderProduct
		for _, product := range order.OrderProducts {
			products = append(products, &OrderProduct{
				ID:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Quantity:    product.Quantity,
			})
		}

		orders = append(orders, &Order{
			ID:         order.ID,
			CreatedAt:  order.CreatedAt.String(),
			TotlaPrice: order.Price,
			Products:   products,
		})
	}

	return orders, nil
}
