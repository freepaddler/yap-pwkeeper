package grpcapi

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"yap-pwkeeper/internal/app/server/aaa"
	"yap-pwkeeper/internal/pkg/grpc/proto"
	"yap-pwkeeper/internal/pkg/jwtToken"
	"yap-pwkeeper/internal/pkg/logger"
	"yap-pwkeeper/internal/pkg/models"
)

type AAA interface {
	Register(ctx context.Context, cred models.UserCredentials) (string, error)
	Login(ctx context.Context, cred models.UserCredentials) (string, error)
	Refresh(ctx context.Context, token string) (string, error)
	Validate(ctx context.Context, token string) bool
}

type AuthHandlers struct {
	proto.UnimplementedAuthServer
	auth AAA
}

func NewAuthHandlers(auth AAA) *AuthHandlers {
	return &AuthHandlers{auth: auth}
}

func (a AuthHandlers) Register(ctx context.Context, in *proto.LoginCredentials) (*proto.Token, error) {
	log := logger.Log().WithCtxRequestId(ctx).With("login", in.Login)
	log.Debug("user registration request")
	response := &proto.Token{}
	credentials := models.UserCredentials{
		Login:    in.Login,
		Password: in.Password,
	}
	token, err := a.auth.Register(ctx, credentials)
	switch {
	case errors.Is(aaa.ErrDuplicate, err):
		return response, status.Error(codes.AlreadyExists, aaa.ErrDuplicate.Error())
	case err != nil:
		return response, status.Error(codes.Internal, "server error")
	}
	response.Token = token
	return response, nil
}

func (a AuthHandlers) Login(ctx context.Context, in *proto.LoginCredentials) (*proto.Token, error) {
	log := logger.Log().WithCtxRequestId(ctx).With("login", in.Login)
	log.Debug("user login request")
	response := &proto.Token{}
	credentials := models.UserCredentials{
		Login:    in.Login,
		Password: in.Password,
	}
	token, err := a.auth.Login(ctx, credentials)
	switch {
	case errors.Is(aaa.ErrBadAuth, err):
		return response, status.Error(codes.Unauthenticated, aaa.ErrBadAuth.Error())
	case err != nil:
		return response, status.Error(codes.Internal, "server error")
	}
	response.Token = token
	return response, nil
}

func (a AuthHandlers) Refresh(ctx context.Context, in *proto.Token) (*proto.Token, error) {
	log := logger.Log().
		WithCtxRequestId(ctx).
		With("sessionId", jwtToken.GetTokenSession(in.Token), "userId", jwtToken.GetTokenSubject(in.Token))
	log.Debug("token refresh request")
	response := &proto.Token{}
	token, err := a.auth.Refresh(ctx, in.Token)
	switch {
	case errors.Is(aaa.ErrBadAuth, err):
		return response, status.Error(codes.Unauthenticated, aaa.ErrBadAuth.Error())
	case err != nil:
		return response, status.Error(codes.Internal, "server error")
	}
	response.Token = token
	return response, nil
}
