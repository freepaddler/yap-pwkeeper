package interceptors

import (
	"context"
	"regexp"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"yap-pwkeeper/internal/pkg/jwtToken"
	"yap-pwkeeper/internal/pkg/logger"
)

const authHeader = "Bearer"

func AuthUnaryServer(valid func(context.Context, string) bool, applyTo ...string) func(ctx context.Context, req interface{}, ServerInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, ServerInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		apply := make(map[string]bool)
		for _, a := range applyTo {
			apply[a] = true
		}
		if !applicable(apply, ServerInfo.FullMethod) {
			logger.Log().WithCtxRequestId(ctx).Debug("skp auth: method does not match")
			return handler(ctx, req)
		}
		logger.Log().WithCtxRequestId(ctx).Debug("check auth")
		// try to get token from request
		var token string
		md, ok := metadata.FromIncomingContext(ctx)
		if ok && md.Len() > 0 {
			h := md.Get(authHeader)
			if len(h) > 0 && len(h[0]) > 0 {
				token = h[0]
			}
		}
		if !valid(ctx, token) {
			return nil, status.Error(codes.Unauthenticated, "")
		}
		ctx = logger.WithUserId(ctx, jwtToken.GetTokenSubject(token))
		return handler(ctx, req)
	}
}

func applicable(apply map[string]bool, si string) bool {
	if len(apply) == 0 {
		return true
	}
	rx, err := regexp.Compile(`.*\.((.*/).*)$`)
	if err != nil {
		return false
	}
	rxMatch := rx.FindAllStringSubmatch(si, -1)
	if rxMatch == nil || len(rxMatch[0]) < 3 {
		return false
	}
	method, service := rxMatch[0][1], rxMatch[0][2]
	if apply[method] || apply[service] {
		return true
	}
	return false
}
