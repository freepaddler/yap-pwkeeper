package aaa

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"yap-pwkeeper/internal/pkg/logger"
	"yap-pwkeeper/internal/pkg/models"
	"yap-pwkeeper/pkg/jwtToken"
)

var (
	ErrNotFound  = errors.New("user not found")
	ErrDuplicate = errors.New("user already exists")
	ErrBadAuth   = errors.New("invalid login credentials")
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
		return "", err
	}
	user, err = c.store.AddUser(ctx,
		models.User{
			Login:        cred.Login,
			PasswordHash: string(pwHash),
			Entity: models.Entity{
				CreatedAt:  time.Now(),
				ModifiedAt: time.Now(),
				State:      models.StateActive,
			},
		})
	if err != nil {
		log.Warnf("failed user registration: %s", err.Error())
		return "", err
	}
	log.With("userId", user.Id).Info("success user registration")
	return newSession(ctx, user.Id)
}

func (c *Controller) Login(ctx context.Context, cred models.UserCredentials) (string, error) {
	log := logger.Log().WithCtxRequestId(ctx).With("login", cred.Login)
	log.Debug("new user login")
	user := models.User{}
	user, err := c.store.GetUserByLogin(ctx, cred.Login)
	if err != nil {
		log.Warnf("failed user login : %s", err.Error())
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(cred.Password)); err != nil {
		log.Warnf("failed user login : %s", err.Error())
		return "", ErrBadAuth
	}
	log.With("userId", user.Id).Info("success user login")
	return newSession(ctx, user.Id)
}

func newSession(ctx context.Context, userid string) (string, error) {
	log := logger.Log().WithCtxRequestId(ctx).With("userId", userid)
	log.Debug("new user session")
	token, err := jwtToken.NewToken(userid)
	if err != nil {
		log.Errorf("failed user session: %s", err.Error())
		return "", ErrToken
	}
	log.With("sessionId", jwtToken.GetTokenSession(token)).Info("success user session")
	return token, err
}

func (c *Controller) Refresh(ctx context.Context, token string) (string, error) {
	log := logger.Log().WithCtxRequestId(ctx).
		With("userId", jwtToken.GetTokenSubject(token), "sessionId", jwtToken.GetTokenSession(token))
	log.Debug("session refresh")
	newToken, err := jwtToken.RefreshToken(token)
	if err != nil {
		if errors.Is(jwtToken.ErrInvalid, err) {
			log.Warnf("failed session refresh: %s", err.Error())
			return "", ErrBadAuth
		} else {
			log.Errorf("failed session refresh: %s", err.Error())
		}
		return "", err
	}
	log.Info("success session refresh")
	return newToken, err
}

func (c *Controller) Validate(token string) bool {
	return jwtToken.ValidateToken(token)
}
