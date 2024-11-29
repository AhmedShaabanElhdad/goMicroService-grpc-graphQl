package account

import (
	"context"

	"github.com/AhmedShaabanElhdad/goMicroService-grpc-graphQl/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	// important to close connection in case of error
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	serviceClient := pb.NewAccountServiceClient(conn)
	return &Client{
		conn:    conn,
		service: serviceClient,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	result, err := c.service.PostAccount(ctx, &pb.PostAccountRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return &Account{
		Name: result.Account.Name,
		ID:   result.Account.Id,
	}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	result, err := c.service.GetAccount(ctx, &pb.GetAccountRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return &Account{
		Name: result.Account.Name,
		ID:   result.Account.Id,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip uint64, taken uint64) ([]Account, error) {
	result, err := c.service.GetAccounts(ctx, &pb.GetAccountsRequest{
		Skip: skip,
	})
	if err != nil {
		return nil, err
	}

	accounts := []Account{}
	for _, account := range result.Accounts {
		accounts = append(accounts, Account{
			Name: account.Name,
			ID:   account.Id,
		})
	}

	return accounts, nil
}
