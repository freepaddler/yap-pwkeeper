package aaa

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"yap-pwkeeper/internal/pkg/jwtToken"
	"yap-pwkeeper/internal/pkg/logger"
	"yap-pwkeeper/internal/pkg/models"
)

var (
	ErrDuplicate = errors.New("user already exists")
	ErrBadAuth   = errors.New("invalid auth credentials")
	ErrNotFound  = errors.New("user not found")
	ErrToken     = errors.New("token generation failed")
)

type UserStorage interface {
	AddUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
}

type Controller struct {
	store UserStorage
}

func New(store UserStorage) *Controller {
	c := &Controller{store: store}
	return c
}

func (c *Controller) Register(ctx context.Context, cred models.UserCredentials) (string, error) {
	log := logger.Log().WithCtxRequestId(ctx).With("login", cred.Login)
	log.Debug("new user registration")
	user := models.User{}
	pwHash, err := bcrypt.GenerateFromPassword([]byte(cred.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate passsword hash: %w", err)
	}
	user, err = c.store.AddUser(ctx,
		models.User{
			Login:        cred.Login,
			PasswordHash: string(pwHash),
			State:        models.StateActive,
		})
	if err != nil {
		log.Warnf("user registration failed: %s", err.Error())
		return "", err
	}
	log.With("userId", user.Id).Info("user registration succeeded")
	return newSession(ctx, user.Id)
}

func (c *Controller) Login(ctx context.Context, cred models.UserCredentials) (string, error) {
	log := logger.Log().WithCtxRequestId(ctx).With("login", cred.Login)
	log.Debug("new user login")
	user := models.User{}
	user, err := c.store.GetUserByLogin(ctx, cred.Login)
	if err != nil {
		log.Warnf("user login failed: %s", err.Error())
		if errors.Is(ErrNotFound, err) {
			return "", ErrBadAuth
		}
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(cred.Password)); err != nil {
		log.Warnf("user login failed: %s", err.Error())
		return "", ErrBadAuth
	}
	log.With("userId", user.Id).Info("user login succeeded")
	return newSession(ctx, user.Id)
}

func newSession(ctx context.Context, userid string) (string, error) {
	log := logger.Log().WithCtxRequestId(ctx).With("userId", userid)
	log.Debug("new user session")
	token, err := jwtToken.NewToken(userid)
	if err != nil {
		log.Errorf("user session failed: %s: %s", ErrToken.Error(), err.Error())
		return "", fmt.Errorf("%w: %w", ErrToken, err)
	}
	log.With("sessionId", jwtToken.GetTokenSession(token)).Info("user session succeeded")
	return token, err
}

func (c *Controller) Refresh(ctx context.Context, token string) (string, error) {
	log := logger.Log().WithCtxRequestId(ctx).
		With("userId", jwtToken.GetTokenSubject(token), "sessionId", jwtToken.GetTokenSession(token))
	log.Debug("session refresh")
	newToken, err := jwtToken.RefreshToken(token)
	if err != nil {
		if errors.Is(jwtToken.ErrInvalid, err) {
			log.Warnf("session refresh failed: %s", err.Error())
			return "", ErrBadAuth
		} else {
			log.Errorf("session refresh failed: %s", err.Error())
		}
		return "", err
	}
	log.Info("session refresh succeeded")
	return newToken, err
}

func (c *Controller) Validate(_ context.Context, token string) bool {
	return jwtToken.Valid(token)
}
