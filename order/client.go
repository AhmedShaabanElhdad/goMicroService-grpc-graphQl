package order

import (
	"context"
	"time"

	pb "github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderClient struct {
	service pb.OrderServiceClient
	conn    *grpc.ClientConn
}

func NewClient(orderUrl string) (*OrderClient, error) {
	conn, err := grpc.NewClient(orderUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewOrderServiceClient(conn)

	return &OrderClient{
		service: client,
		conn:    conn,
	}, nil

}

func (client OrderClient) Close() {
	client.conn.Close()
}

func (client OrderClient) CreateOrder(ctx context.Context, products []RequestedOrderProduct, accountId string) (*Order, error) {

	productsRequest := []*pb.PostOrderRequest_OrderProduct{}
	for _, product := range products {
		productsRequest = append(productsRequest, &pb.PostOrderRequest_OrderProduct{
			ProductId: product.ID,
			Quantity:  uint32(product.Quantity),
		})
	}

	order, err := client.service.PostOrder(ctx, &pb.PostOrderRequest{
		AccountId: accountId,
		Products:  productsRequest,
	})
	if err != nil {
		return nil, err
	}
	productsRes := []OrderProduct{}
	for _, productRes := range order.Order.Products {
		productsRes = append(productsRes, OrderProduct{
			ID:          productRes.Id,
			Price:       productRes.Price,
			Name:        productRes.Name,
			Description: productRes.Description,
			Quantity:    int(productRes.Quantity),
		})
	}

	createdAt := time.Time{}
	createdAt.UnmarshalBinary(order.Order.CreatedAt)
	return &Order{
		ID:            order.Order.Id,
		Price:         order.Order.Price,
		AccountID:     order.Order.AccountId,
		CreatedAt:     createdAt,
		OrderProducts: productsRes,
	}, nil
}

func (client OrderClient) GetOrderById(ctx context.Context, id string) (*Order, error) {
	res, err := client.service.GetOrderId(ctx, &pb.GetOrderRequest{Id: id})
	if err != nil {
		return nil, err
	}

	createdAt := time.Time{}
	createdAt.UnmarshalBinary(res.Order.CreatedAt)

	productsRes := []OrderProduct{}
	for _, productRes := range res.Order.Products {
		productsRes = append(productsRes, OrderProduct{
			ID:          productRes.Id,
			Price:       productRes.Price,
			Name:        productRes.Name,
			Description: productRes.Description,
			Quantity:    int(productRes.Quantity),
		})
	}

	return &Order{
		ID:            res.Order.Id,
		Price:         res.Order.Price,
		AccountID:     res.Order.AccountId,
		CreatedAt:     createdAt,
		OrderProducts: productsRes,
	}, nil
}

func (client OrderClient) GetAccountOrder(ctx context.Context, accountId string) ([]Order, error) {
	res, err := client.service.GetAccountOrders(ctx, &pb.GetAccountOrdersRequest{AccountId: accountId})
	if err != nil {
		return nil, err
	}

	var orders []Order
	for _, order := range res.Orders {
		createdAt := time.Time{}
		createdAt.UnmarshalBinary(order.CreatedAt)

		productsRes := []OrderProduct{}
		for _, productRes := range order.Products {
			productsRes = append(productsRes, OrderProduct{
				ID:          productRes.Id,
				Price:       productRes.Price,
				Name:        productRes.Name,
				Description: productRes.Description,
				Quantity:    int(productRes.Quantity),
			})
		}

		orders = append(orders, Order{
			ID:            order.Id,
			Price:         order.Price,
			AccountID:     order.AccountId,
			CreatedAt:     createdAt,
			OrderProducts: productsRes,
		})
	}
	return orders, nil
}

func (client OrderClient) CreateAccountOrders(ctx context.Context, accountId string) ([]*Order, error) {
	res, err := client.service.GetAccountOrders(ctx, &pb.GetAccountOrdersRequest{
		AccountId: accountId,
	})
	if err != nil {
		return nil, err
	}

	orders := []*Order{}

	for _, order := range res.Orders {
		productsRes := []OrderProduct{}
		for _, productRes := range order.Products {

			productsRes = append(productsRes, OrderProduct{
				ID:          productRes.Id,
				Price:       productRes.Price,
				Name:        productRes.Name,
				Description: productRes.Description,
				Quantity:    int(productRes.Quantity),
			})
		}
		createdAt := time.Time{}
		createdAt.UnmarshalBinary(order.CreatedAt)
		orders = append(orders, &Order{
			ID:            order.Id,
			Price:         order.Price,
			AccountID:     order.AccountId,
			CreatedAt:     createdAt,
			OrderProducts: productsRes,
		})
	}
	return orders, nil
}
