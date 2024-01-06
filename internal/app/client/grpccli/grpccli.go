package grpccli

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"yap-pwkeeper/internal/pkg/grpc/proto"
)

type Client struct {
	address     string
	conn        *grpc.ClientConn
	auth        proto.AuthClient
	docs        proto.WalletClient
	authTimeout time.Duration
	docsTimeout time.Duration
}

func New(address string, options ...func(c *Client)) (*Client, error) {
	cli := &Client{
		address:     address,
		authTimeout: 5 * time.Second,
		docsTimeout: 30 * time.Second,
	}
	for _, opt := range options {
		opt(cli)
	}
	var err error
	cli.conn, err = grpc.Dial(
		cli.address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	cli.auth = proto.NewAuthClient(cli.conn)
	cli.docs = proto.NewWalletClient(cli.conn)
	return cli, err
}

func WithTimeouts(auth, docs time.Duration) func(c *Client) {
	return func(c *Client) {
		c.authTimeout = auth
		c.docsTimeout = docs
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Register(login string, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.authTimeout)
	defer cancel()
	cred := &proto.LoginCredentials{
		Login:    login,
		Password: password,
	}
	token := &proto.Token{}
	token, err := c.auth.Register(ctx, cred)
	return token.GetToken(), err
}

func (c *Client) Login(login string, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.authTimeout)
	defer cancel()
	cred := &proto.LoginCredentials{
		Login:    login,
		Password: password,
	}
	token := &proto.Token{}
	token, err := c.auth.Login(ctx, cred)
	return token.GetToken(), err
}

func (c *Client) RefreshToken(token string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.authTimeout)
	defer cancel()
	oldToken := &proto.Token{
		Token: token,
	}
	newToken := &proto.Token{}
	newToken, err := c.auth.Refresh(ctx, oldToken)
	return newToken.GetToken(), err
}
