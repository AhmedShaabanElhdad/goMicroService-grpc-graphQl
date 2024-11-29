package catalog

import (
	"context"
	"fmt"
	"net"

	pb "github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type catalogGrpcServer struct {
	pb.UnimplementedProductServiceServer
	service Service
}

func ListenAndServeGrpc(service Service, port int) error {
	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil
	}
	server := grpc.NewServer()
	pb.RegisterProductServiceServer(server, &catalogGrpcServer{
		service: service,
	})

	reflection.Register(server)
	return server.Serve(list)
}

func (server catalogGrpcServer) Close() {

}
func (server catalogGrpcServer) PostCatalog(ctx context.Context, r *pb.PostProductRequest) error {
	return server.service.PostCatalog(ctx, r.Name, r.Description, r.Price)
}

func (server catalogGrpcServer) GetProducts(ctx context.Context, r *pb.GetProductsRequest) (*pb.ProductsResponse, error) {
	products, err := server.service.GetProducts(ctx,
		uint64(r.Taken),
		uint64(r.Skip),
	)
	if err != nil {
		return nil, err
	}

	res := []*pb.Product{}
	for _, product := range products {
		res = append(res, &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	return &pb.ProductsResponse{
		Products: res,
	}, nil
}

func (server catalogGrpcServer) GetProductsByID(ctx context.Context, r *pb.GetProductsByIDRequest) (*pb.ProductResponse, error) {
	product, err := server.service.GetProductsByID(ctx,
		r.Id,
	)
	if err != nil {
		return nil, err
	}
	return &pb.ProductResponse{
		Product: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		},
	}, nil
}

func (server catalogGrpcServer) GetProductsWithIds(ctx context.Context, r *pb.GetProductsWithIdsRequest) (*pb.ProductsResponse, error) {
	products, err := server.service.GetProductsWithIds(ctx, r.Id)
	if err != nil {
		return nil, err
	}

	res := []*pb.Product{}
	for _, product := range products {
		res = append(res, &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	return &pb.ProductsResponse{
		Products: res,
	}, nil

}

func (server catalogGrpcServer) SearchProduct(ctx context.Context, r *pb.SearchProductRequest) (*pb.ProductsResponse, error) {
	products, err := server.service.SearchProduct(ctx, r.Query, uint64(r.Taken), uint64(r.Skip))
	if err != nil {
		return nil, err
	}

	res := []*pb.Product{}
	for _, product := range products {
		res = append(res, &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	return &pb.ProductsResponse{
		Products: res,
	}, nil
}
