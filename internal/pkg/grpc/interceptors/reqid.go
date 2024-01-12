package interceptors

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"yap-pwkeeper/internal/pkg/logger"
)

const requestIdHeader = "request-id"

// ReqIdStreamServer assigns unique identification to each streaming request
func ReqIdStreamServer(srv interface{}, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx, mdOut := setReqId(stream.Context())
	if err := grpc.SetHeader(ctx, mdOut); err != nil {
		logger.Log().WithErr(err).WithCtxRequestId(ctx).Error("failed to set grpc headers")
	}
	return handler(srv, &serverStreamWrapped{stream, ctx})
}

// ReqIdUnaryServer is the same as ReqIdStreamServer, but for unary requests
func ReqIdUnaryServer(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	ctx, mdOut := setReqId(ctx)
	if err := grpc.SetHeader(ctx, mdOut); err != nil {
		logger.Log().WithErr(err).Error("failed to set grpc headers")
	}
	return handler(ctx, req)
}

func setReqId(ctx context.Context) (context.Context, metadata.MD) {
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
	return ctx, mdOut
}
