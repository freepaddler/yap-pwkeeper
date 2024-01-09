package grpcapi

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"yap-pwkeeper/internal/app/server/documents"
	"yap-pwkeeper/internal/pkg/grpc/interceptors"
	pb "yap-pwkeeper/internal/pkg/grpc/proto"
	"yap-pwkeeper/internal/pkg/logger"
)

type GRCPServer struct {
	Server             *grpc.Server
	address            string
	auth               *AuthHandlers
	wallet             *DocsHandlers
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
}

func (gs *GRCPServer) GetAddress() string {
	return gs.address
}

func New(opts ...func(gs *GRCPServer)) *GRCPServer {
	gs := new(GRCPServer)
	// setup logging
	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.FinishCall),
		logging.WithFieldsFromContext(interceptors.LoggingFields),
	}
	gs.unaryInterceptors = append(
		gs.unaryInterceptors,
		interceptors.ReqIdUnaryServer,
		logging.UnaryServerInterceptor(interceptors.ZapLogger(logger.Log().Desugar()), logOpts...),
	)
	gs.streamInterceptors = append(
		gs.streamInterceptors,
		interceptors.ReqIdStreamServer,
		logging.StreamServerInterceptor(interceptors.ZapLogger(logger.Log().Desugar()), logOpts...),
	)
	for _, o := range opts {
		o(gs)
	}
	gs.Server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(gs.unaryInterceptors...),
		grpc.ChainStreamInterceptor(gs.streamInterceptors...),
	)
	pb.RegisterAuthServer(gs.Server, gs.auth)
	pb.RegisterDocsServer(gs.Server, gs.wallet)
	return gs
}

func WithAddress(s string) func(gs *GRCPServer) {
	return func(gs *GRCPServer) {
		gs.address = s
	}
}

func WithAuthHandlers(h *AuthHandlers) func(gs *GRCPServer) {
	return func(gs *GRCPServer) {
		gs.auth = h
	}
}

func WithDocsHandlers(h *DocsHandlers) func(gs *GRCPServer) {
	return func(gs *GRCPServer) {
		gs.wallet = h
	}
}

func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) func(server *GRCPServer) {
	return func(gs *GRCPServer) {
		// order makes sense
		gs.unaryInterceptors = append(gs.unaryInterceptors, interceptors...)
	}
}

func WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) func(server *GRCPServer) {
	return func(gs *GRCPServer) {
		// order makes sense
		gs.streamInterceptors = append(gs.streamInterceptors, interceptors...)
	}
}

// respErr returns grpc error response
func respErr(ctx context.Context, err error) error {
	switch {
	case errors.Is(documents.ErrBadRequest, err):
		return status.Error(codes.InvalidArgument, documents.ErrBadRequest.Error())
	case errors.Is(documents.ErrChanged, err):
		return status.Error(codes.FailedPrecondition, documents.ErrChanged.Error())
	case errors.Is(documents.ErrDeleted, err):
		return status.Error(codes.FailedPrecondition, documents.ErrDeleted.Error())
	case errors.Is(documents.ErrNotFound, err):
		return status.Error(codes.NotFound, documents.ErrNotFound.Error())
	default:
		logger.Log().WithErr(err).WithCtxRequestId(ctx).Error("server error")
		return status.Error(codes.Internal, "server error")
	}
}
