package grpccli

import (
	"context"
	"errors"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"yap-pwkeeper/internal/pkg/grpc/proto"
	"yap-pwkeeper/internal/pkg/jwtToken"
)

// Register registers new client on server
func (c *Client) Register(login, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.authTimeout)
	defer cancel()
	cred := &proto.LoginCredentials{
		Login:    login,
		Password: password,
	}
	token := &proto.Token{}
	token, err := c.auth.Register(ctx, cred)
	if err != nil {
		return parseErr(err)
	}
	c.setToken(token.GetToken())
	return nil
}

// Login logins to server and starts token update routine
// It is safe to call login multiple times.
// After each successful attempt new token will be used.
func (c *Client) Login(login, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.authTimeout)
	defer cancel()
	cred := &proto.LoginCredentials{
		Login:    login,
		Password: password,
	}
	token := &proto.Token{}
	token, err := c.auth.Login(ctx, cred)
	if err != nil {
		return parseErr(err)
	}
	c.setToken(token.GetToken())
	return nil
}

// Logout flushes token from memory.
// This will also terminate refresh routine.
// Safe to be called multiple times.
func (c *Client) Logout() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = ""
}

// getToken safely returns current token
// this is the ONLY way to get curren token
func (c *Client) getToken() string {
	c.mu.RLock()
	token := c.token
	c.mu.RUnlock()
	return token
}

// setToken safely updates token and launches token update routine
// it controls that only one update process is running
// passing
func (c *Client) setToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	close(c.ch)
	c.token = token
	c.ch = make(chan struct{})
	log.Println("got new token")
	go c.tokenUpdater(token)
}

// refreshToken implements server method to refresh token
func (c *Client) refreshToken(token string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.authTimeout)
	defer cancel()
	oldToken := &proto.Token{
		Token: token,
	}
	newToken := &proto.Token{}
	newToken, err := c.auth.Refresh(ctx, oldToken)
	return newToken.GetToken(), err
}

// tokenUpdater starts token update procedure when tokenTimeUntilExpire left for the token.
// Update attempts scheduled every tokenRefreshRetryInterval
// tokenUpdater exits on the following conditions:
//   - token is an empty string
//   - token is not valid (claims can't be parsed) or does not have exp claim
//   - server rejected token update with Unauthenticated response
func (c *Client) tokenUpdater(token string) {
	log.Println("token refresh routine started")
	var err error
	defer func() {
		if err != nil {
			log.Printf("token refresh routine terminated: %s", err.Error())
		} else {
			log.Println("token refresh routine stopped")
		}
	}()
	if token == "" {
		err = errors.New("empty token")
		return
	}
	expire, err := jwtToken.GetTokenExpire(token)
	if err != nil {
		err = errors.New("unable to parse exp from token")
		return
	}
	var newToken string
	// token refresh routine
	select {
	case <-c.ch:
		return
	// start refresh cycle when timeToExpire left
	case <-time.After(time.Until(expire) - c.tokenTimeUntilExpire):
		for {
			newToken, err = c.refreshToken(token)
			st, _ := status.FromError(err)
			switch st.Code() {
			case codes.OK:
				c.setToken(newToken)
				return
			case codes.Unauthenticated:
				err = errors.New("token rejected")
				return
			default:
				log.Printf("token refresh failed: %s", err.Error())
			}
			err = nil
			select {
			case <-c.ch:
				return
			case <-time.After(c.tokenRefreshRetryInterval):
				if expire.Before(time.Now()) {
					log.Println("token is already expired, but why not to try... ")
				}
			}
		}
	}
}
