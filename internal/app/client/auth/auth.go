package auth

import (
	"errors"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"yap-pwkeeper/internal/pkg/jwtToken"
)

type AuthServer interface {
	register(login string, password string) (string, error)
	Login(login string, password string) (string, error)
	RefreshToken(token string) (string, error)
}

var (
	ErrExpired  = errors.New("token expired")
	ErrRejected = errors.New("token rejected")
)

type Controller struct {
	refreshBeforeExpire  time.Duration
	refreshRetryInterval time.Duration
	server               AuthServer
	token                string
	ch                   chan struct{}
	mu                   sync.RWMutex
}

func New(server AuthServer) *Controller {
	return &Controller{
		server:               server,
		ch:                   make(chan struct{}),
		refreshBeforeExpire:  2 * time.Minute,
		refreshRetryInterval: 5 * time.Second,
	}
}

func (c *Controller) GetToken() string {
	c.mu.RLock()
	token := c.token
	c.mu.RUnlock()
	return token
}

func (c *Controller) Register(login, password string) error {
	token, err := c.server.register(login, password)
	if err != nil {
		return err
	}
	c.setToken(token)
	return nil
}

func (c *Controller) Login(login, password string) error {
	token, err := c.server.Login(login, password)
	if err != nil {
		return err
	}
	c.setToken(token)
	return nil
}

func (c *Controller) setToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	close(c.ch)
	c.token = token
	c.ch = make(chan struct{})
	log.Println("got new token")
	go c.refreshToken(token)
}

func (c *Controller) refreshToken(token string) {
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
			newToken, err = c.server.RefreshToken(token)
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
