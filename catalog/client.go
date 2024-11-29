package catalog

import (
	"context"

	"github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductClient struct {
	conn    *grpc.ClientConn
	service pb.ProductServiceClient
}

func NewClient(url string) (*ProductClient, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	service := pb.NewProductServiceClient(conn)
	client := &ProductClient{
		conn:    conn,
		service: service,
	}
	return client, nil
}

func (c ProductClient) Close() {
	c.conn.Close()
}
func (c ProductClient) PostCatalog(ctx context.Context, catalog Catalog) (*Catalog, error) {
	res, err := c.service.PostProduct(ctx, &pb.PostProductRequest{
		Name:        catalog.Name,
		Description: catalog.Description,
		Price:       catalog.Price,
	})
	if err != nil {
		return nil, err
	}

	return &Catalog{
		ID:          res.Product.Id,
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Price:       res.Product.Price,
	}, err
}
func (c ProductClient) GetProducts(ctx context.Context, taken uint64, skip uint64, query string) ([]Catalog, error) {

	var (
		res *pb.ProductsResponse
		err error
	)
	if query != "" {
		res, err = c.service.GetProducts(ctx, &pb.GetProductsRequest{
			Taken: int64(taken),
			Skip:  int64(skip),
		})
	} else {
		res, err = c.service.SearchProduct(ctx, &pb.SearchProductRequest{
			Taken: int64(taken),
			Skip:  int64(skip),
			Query: query,
		})
	}

	if err != nil {
		return nil, err
	}
	catalogs := []Catalog{}

	for _, product := range res.Products {
		catalogs = append(catalogs, Catalog{
			ID:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       float64(product.Price),
		})
	}

	return catalogs, nil
}
func (c ProductClient) GetProductsByID(ctx context.Context, id string) (*Catalog, error) {
	res, err := c.service.GetProductsByID(ctx, &pb.GetProductsByIDRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}
	catalog := &Catalog{
		ID:          res.Product.Id,
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Price:       float64(res.Product.Price),
	}

	return catalog, err
}
func (c ProductClient) GetProductsWithIds(ctx context.Context, ids []string) ([]Catalog, error) {

	res, err := c.service.GetProductsWithIds(ctx, &pb.GetProductsWithIdsRequest{
		Id: ids,
	})

	if err != nil {
		return nil, err
	}
	catalogs := []Catalog{}

	for _, product := range res.Products {
		catalogs = append(catalogs, Catalog{
			ID:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       float64(product.Price),
		})
	}

	return catalogs, nil
}

func (c ProductClient) SearchProduct(ctx context.Context, query string, taken uint64, skip uint64) ([]Catalog, error) {
	res, err := c.service.SearchProduct(ctx, &pb.SearchProductRequest{
		Query: query,
		Taken: int64(taken),
		Skip:  int64(skip),
	})

	if err != nil {
		return nil, err
	}
	catalogs := []Catalog{}

	for _, product := range res.Products {
		catalogs = append(catalogs, Catalog{
			ID:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       float64(product.Price),
		})
	}

	return catalogs, nil
}
