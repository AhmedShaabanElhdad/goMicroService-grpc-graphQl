package main

import (
	"context"
	"log"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(3*time.Second))
	defer cancel()

	if id != nil {
		res, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*Account{
			{
				ID:   res.ID,
				Name: res.Name,
			},
		}, nil
	}

	skip, taken := uint64(0), uint64(0)
	if pagination != nil {
		skip, taken = pagination.bounds()
	}

	accountList, err := r.server.accountClient.GetAccounts(ctx, skip, taken)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var accounts []*Account
	for _, account := range accountList {
		accounts = append(accounts, &Account{
			ID:   account.ID,
			Name: account.Name,
		})
	}

	return accounts, nil
}

func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(3*time.Second))
	defer cancel()

	if id != nil {
		res, err := r.server.catalogClient.GetProductsByID(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*Product{
			{
				ID:          res.ID,
				Name:        res.Name,
				Description: res.Description,
				Price:       res.Price,
			},
		}, nil
	}

	skip, taken := uint64(0), uint64(0)
	if pagination != nil {
		skip, taken = pagination.bounds()
	}

	q := ""
	if query != nil {
		q = *query
	}

	productList, err := r.server.catalogClient.GetProducts(ctx, skip, taken, q)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var products []*Product
	for _, product := range productList {
		products = append(products, &Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	return products, nil
}

func (p PaginationInput) bounds() (uint64, uint64) {

	skipValue := uint64(0)
	takenValue := uint64(0)

	if p.Skip != nil {
		skipValue = uint64(*p.Skip)
	}
	if p.Take != nil {
		takenValue = uint64(*p.Skip)
	}

	return skipValue, takenValue
}
