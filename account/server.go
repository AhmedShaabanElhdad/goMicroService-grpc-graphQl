package account

import (
	"context"
	"fmt"
	"net"

	pb "github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	service Service
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterAccountServiceServer(serv, &grpcServer{
		service: s,
	})

	//todo search for this part
	reflection.Register(serv)
	return serv.Serve(lis)
}

// func NewServer() *Server {
// 	return &Server{
// 		service: &accounService{},
// 	}
// }

func (s *grpcServer) GetAccount(ctx context.Context, r *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	a, err := s.service.GetAccount(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetAccountResponse{
		Account: &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		},
	}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, r *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	res, err := s.service.GetAccounts(ctx, r.Skip, r.Take)
	if err != nil {
		return nil, err
	}

	accounts := []*pb.Account{}
	for _, account := range res {
		accounts = append(accounts, &pb.Account{
			Id:   account.ID,
			Name: account.Name,
		})

	}

	return &pb.GetAccountsResponse{
		Accounts: accounts,
	}, nil

}

func (s *grpcServer) PostAccount(ctx context.Context, r *pb.PostAccountRequest) (*pb.PostAccountResponse, error) {
	a, err := s.service.PostAccoun(ctx, r.Name)
	if err != nil {
		return nil, err
	}
	return &pb.PostAccountResponse{
		Account: &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		},
	}, nil
}
