package interceptors

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"yap-pwkeeper/internal/pkg/logger"
)

const requestIdHeader = "request-id"

func ReqIdUnaryServer(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	var requestId string

	// try to get requestId from request
	md, ok := metadata.FromIncomingContext(ctx)
	if ok && md.Len() > 0 {
		h := md.Get(requestIdHeader)
		if len(h) > 0 && len(h[0]) > 0 {
			requestId = h[0]
		}
	}

	// create metadata if not exists
	if !ok {
		md = metadata.New(map[string]string{})
	}

	// no requestId in header
	if requestId == "" {
		requestId = uuid.NewString()
		// put requestId in context
		md.Append(requestIdHeader, requestId)
		ctx = metadata.NewIncomingContext(ctx, md)
	}

	// put requestId in context
	ctx = logger.WithRequestId(ctx, requestId)

	// append same requestId to response
	mdOut := metadata.New(map[string]string{requestIdHeader: requestId})
	if err := grpc.SetHeader(ctx, mdOut); err != nil {
		logger.Log().WithErr(err).Error("failed to set grpc headers")
	}

	return handler(ctx, req)
}
