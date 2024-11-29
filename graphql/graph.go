package main

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl/account"
	"github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl/catalog"
	"github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl/order"
)

type Server struct {
	accountClient *account.Client
	catalogClient *catalog.ProductClient
	orderClient   *order.OrderClient
}

func NewGraphQlServer(accountUrl, catalogUrl, orderUrl string) (*Server, error) {
	accountClient, err := account.NewClient(accountUrl)
	if err != nil {
		return nil, err
	}
	catalogClient, err := catalog.NewClient(catalogUrl)
	if err != nil {
		accountClient.Close()
		return nil, err
	}

	orderClient, err := order.NewClient(orderUrl)
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return nil, err
	}

	return &Server{
		accountClient,
		catalogClient,
		orderClient,
	}, nil
}

func (s *Server) Mutation() MutationResolver {
	return &mutationResolver{
		server: s,
	}
}

func (s *Server) Query() QueryResolver {
	return &queryResolver{
		server: s,
	}
}

func (s *Server) Account() AccountResolver {
	return &accountResolver{
		server: s,
	}
}

func (s *Server) toExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}
