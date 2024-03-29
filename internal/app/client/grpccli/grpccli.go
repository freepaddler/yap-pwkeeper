// Package grpccli is a gRCP client to the server.
// It handles automated token refresh routine after successful login.
// When any of methods return ErrAuthFail error, this means that
// auth session with server is closed and client should Login again.
// Login and Logout methods are safe to be used multiple times.
// After Close is called, this client instance becomes unusable.
package grpccli

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	"yap-pwkeeper/internal/pkg/grpc/proto"
)

type Client struct {
	address                   string                           // server address
	tlsCredentials            credentials.TransportCredentials // tls setup
	conn                      *grpc.ClientConn                 // server connection
	auth                      proto.AuthClient                 // server auth service
	docs                      proto.DocsClient                 // server documents service
	authTimeout               time.Duration                    // timeout for auth service
	dataTimeout               time.Duration                    // timeout for documents service
	tokenTimeUntilExpire      time.Duration                    // time left until token expired
	tokenRefreshRetryInterval time.Duration                    // token refresh retry
	token                     string
	ch                        chan struct{} // token refresher control chan
	mu                        sync.RWMutex  // token and chan mutex
}

// New is a gRPCClient constructor. Address is server connection endpoint `host:port`.
func New(address string, options ...func(c *Client)) (*Client, error) {
	cli := &Client{
		address:                   address,
		authTimeout:               5 * time.Second,
		dataTimeout:               30 * time.Second,
		ch:                        make(chan struct{}),
		tokenTimeUntilExpire:      2 * time.Minute,
		tokenRefreshRetryInterval: 5 * time.Second,
	}
	for _, opt := range options {
		opt(cli)
	}
	var err error
	cli.conn, err = grpc.Dial(
		cli.address,
		grpc.WithTransportCredentials(cli.tlsCredentials),
	)
	cli.auth = proto.NewAuthClient(cli.conn)
	cli.docs = proto.NewDocsClient(cli.conn)
	return cli, err
}

// WithTimeouts sets server requests timeouts.
// authTimeout - for auth service (5 seconds).
// dataTimeout - for documents service (30 seconds).
func WithTimeouts(authTimeout, dataTimeout time.Duration) func(c *Client) {
	return func(c *Client) {
		c.authTimeout = authTimeout
		c.dataTimeout = dataTimeout
	}
}

// WithTokenRefresh sets server token refresh options.
// Refresh starts when tokenTimeUntilExpire time left until expiration (2 minutes).
// tokenRefreshRetryInterval defines time between retries if refresh failed
func WithTokenRefresh(tokenTimeUntilExpire, tokenRefreshRetryInterval time.Duration) func(c *Client) {
	return func(c *Client) {
		c.tokenTimeUntilExpire = tokenTimeUntilExpire
		c.tokenRefreshRetryInterval = tokenRefreshRetryInterval
	}
}

// WithTransportCredentials sets up connection transport security
func WithTransportCredentials(cred credentials.TransportCredentials) func(c *Client) {
	return func(c *Client) {
		c.tlsCredentials = cred
	}
}

// Close closes server connection.
// Instance can't be reused after Close is called.
func (c *Client) Close() error {
	c.Logout()
	return c.conn.Close()
}

var (
	// ErrAuthFail error indicates that client is no more authenticated on server.
	// User Login method to authenticate again
	ErrAuthFail    = errors.New("authorization failed")
	ErrUnavailable = errors.New("unable to connect server")
)

// parseErr returns parsed gRPCErrors
func parseErr(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	switch st.Code() {
	case codes.Unauthenticated:
		return ErrAuthFail
	case codes.Unavailable:
		return fmt.Errorf("%w: %w", ErrUnavailable, err)
	default:
		if st.Message() == "" {
			return errors.New(st.Code().String())
		}
		return errors.New(st.Message())
	}
}

// LoadCACertificate reads CA certificate from file and returns secure config for gRPC client
// insecure flag disables verification of server certificate
func LoadCACertificate(caFile string, insecure bool) (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	caCertificate, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCertificate) {
		return nil, fmt.Errorf("failed to add CA certificate: %w", err)
	}

	config := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: insecure,
	}

	return credentials.NewTLS(config), nil
}
