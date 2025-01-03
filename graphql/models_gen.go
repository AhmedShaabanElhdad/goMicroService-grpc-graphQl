// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package main

type AccountInput struct {
	Name string `json:"name"`
}

type Mutation struct {
}

type NewUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Order struct {
	ID         string          `json:"id"`
	CreatedAt  string          `json:"createdAt"`
	TotlaPrice float64         `json:"totlaPrice"`
	Products   []*OrderProduct `json:"products"`
}

type OrderInput struct {
	AccountID string               `json:"accountId"`
	Products  []*OrderProductInput `json:"products"`
}

type OrderProduct struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type OrderProductInput struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type PaginationInput struct {
	Skip *int `json:"skip,omitempty"`
	Take *int `json:"take,omitempty"`
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type ProductInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type Query struct {
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
