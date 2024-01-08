package grpccli

import (
	"errors"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"yap-pwkeeper/internal/pkg/grpc/proto"
)

type Client struct {
	address              string
	conn                 *grpc.ClientConn
	auth                 proto.AuthClient
	docs                 proto.WalletClient
	authTimeout          time.Duration
	docsTimeout          time.Duration
	refreshBeforeExpire  time.Duration
	refreshRetryInterval time.Duration
	token                string
	ch                   chan struct{}
	mu                   sync.RWMutex
}

func New(address string, options ...func(c *Client)) (*Client, error) {
	cli := &Client{
		address:              address,
		authTimeout:          5 * time.Second,
		docsTimeout:          30 * time.Second,
		ch:                   make(chan struct{}),
		refreshBeforeExpire:  2 * time.Minute,
		refreshRetryInterval: 5 * time.Second,
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

var (
	ErrAuthFail    = errors.New("authorization failed, new login required")
	ErrUnavailable = errors.New("unable to connect server")
)

func parseErr(err error) error {
	if err == nil {
		return nil
	}
	st, _ := status.FromError(err)
	switch st.Code() {
	case codes.Unauthenticated:
		return ErrAuthFail
	case codes.Unavailable:
		return ErrUnavailable
	default:
		if st.Message() == "" {
			return errors.New(st.Code().String())
		}
		return errors.New(st.Message())
	}
}
