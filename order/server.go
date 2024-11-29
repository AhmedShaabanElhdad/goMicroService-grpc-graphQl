package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl/account"
	"github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl/catalog"
	pb "github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	service       Service
	catalogClient catalog.ProductClient
	accountClient account.Client
}

func ListenAndServe(s Service, port int, catalogClient catalog.ProductClient, accountClient account.Client) error {

	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("Error")
	}
	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, NewOrderServer(
		s, catalogClient, accountClient,
	))

	reflection.Register(serv)
	return serv.Serve(list)
}

func NewOrderServer(s Service, catalogClient catalog.ProductClient, accountClient account.Client) *OrderServer {
	return &OrderServer{
		service:       s,
		accountClient: accountClient,
		catalogClient: catalogClient,
	}
}

func (s *OrderServer) GetAccountOrders(ctx context.Context, request *pb.GetAccountOrdersRequest) (*pb.GetAccountOrdersResponse, error) {

	_, err := s.accountClient.GetAccount(ctx, request.AccountId)
	if err != nil {
		log.Printf("Error Getting Account: %v", err)
		return nil, errors.New("account not Found")
	}

	orders, err := s.service.GetAccountOrders(ctx, request.AccountId)
	if err != nil {
		return nil, err
	}

	var ordersMessage []*pb.Order
	for _, order := range orders {

		products := []*pb.Order_OrderProduct{}

		for _, orderProduct := range order.OrderProducts {
			products = append(products, &pb.Order_OrderProduct{
				Id:          orderProduct.ID,
				Name:        orderProduct.Name,
				Description: orderProduct.Description,
				Price:       orderProduct.Price,
				Quantity:    uint32(orderProduct.Quantity),
			})
		}

		createAtInBytes, _ := order.CreatedAt.MarshalBinary()
		orderResponse := pb.Order{
			Id:        order.ID,
			Price:     order.Price,
			AccountId: order.AccountID,
			CreatedAt: createAtInBytes,
			Products:  products,
		}

		ordersMessage = append(ordersMessage, &orderResponse)
	}
	return &pb.GetAccountOrdersResponse{
		Orders: ordersMessage,
	}, nil
}
func (s *OrderServer) GetOrderId(ctx context.Context, request *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	return nil, nil
}
func (s *OrderServer) PostOrder(ctx context.Context, request *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	account, err := s.accountClient.GetAccount(ctx, request.AccountId)
	if err != nil {
		log.Printf("Error Getting Account: %v", err)
		return nil, errors.New("account not Found")
	}

	productIds := []string{}
	for _, product := range request.Products {
		productIds = append(productIds, product.ProductId)
	}

	products, err := s.catalogClient.GetProductsWithIds(ctx, productIds)
	if err != nil {
		log.Printf("Error Getting Account: %v", err)
		return nil, errors.New("there is Wrong Product")
	}

	orderProducts := []OrderProduct{}
	for _, product := range products {
		// todo check quantity
		orderProduct := OrderProduct{
			ID:          product.ID,
			Price:       product.Price,
			Name:        product.Name,
			Description: product.Description,
		}

		// better to find way to use map
		for _, reqProduct := range request.Products {
			if reqProduct.ProductId == product.ID {
				orderProduct.Quantity = int(reqProduct.Quantity)
				break
			}
		}
		if orderProduct.Quantity != 0 {
			orderProducts = append(orderProducts, orderProduct)
		}

	}

	createdOrder, err := s.service.PostOrder(ctx, account.ID, orderProducts)
	if err != nil {
		log.Printf("Error Getting Account: %v", err)
		return nil, err
	}

	order_OrderProduct := []*pb.Order_OrderProduct{}
	for _, createdOrderProduct := range createdOrder.OrderProducts {
		order_OrderProduct = append(order_OrderProduct, &pb.Order_OrderProduct{
			Id:          createdOrderProduct.ID,
			Name:        createdOrderProduct.Name,
			Description: createdOrderProduct.Description,
			Price:       createdOrder.Price,
			Quantity:    uint32(createdOrderProduct.Quantity),
		})
	}

	createdAtBinary, _ := createdOrder.CreatedAt.MarshalBinary()
	return &pb.PostOrderResponse{Order: &pb.Order{
		Id:        createdOrder.ID,
		Price:     createdOrder.Price,
		AccountId: createdOrder.AccountID,
		CreatedAt: createdAtBinary,
		Products:  order_OrderProduct,
	}}, nil
}
