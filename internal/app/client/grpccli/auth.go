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

var (
	ErrExpired  = errors.New("token expired")
	ErrRejected = errors.New("token rejected")
)

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

func (c *Client) getToken() string {
	c.mu.RLock()
	token := c.token
	c.mu.RUnlock()
	return token
}

func (c *Client) setToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	close(c.ch)
	c.token = token
	c.ch = make(chan struct{})
	log.Println("got new token")
	go c.tokenUpdater(token)
}

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

func (c *Client) tokenUpdater(token string) {
	log.Println("token refresh routine started")
	var err error
	var newToken string
	defer func() {
		if err != nil {
			log.Printf("token refresh routine terminated: %s", err.Error())
		} else {
			log.Println("token refresh routine stopped")
		}
	}()
	expire, err := jwtToken.GetTokenExpire(token)
	if err != nil {
		return
	}
	select {
	case <-c.ch:
		return
	case <-time.After(time.Until(expire) - c.refreshBeforeExpire):
		for {
			newToken, err = c.refreshToken(token)
			st, _ := status.FromError(err)
			switch st.Code() {
			case codes.OK:
				c.setToken(newToken)
				return
			case codes.Unauthenticated:
				err = ErrRejected
				return
			default:
				log.Printf("token refresh failed: %s", err.Error())
			}
			err = nil
			select {
			case <-c.ch:
				return
			case <-time.After(time.Until(expire) + time.Second):
				err = ErrExpired
				return
			case <-time.After(c.refreshRetryInterval):
			}
		}
	}
}
