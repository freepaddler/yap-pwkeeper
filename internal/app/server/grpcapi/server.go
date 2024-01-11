// Package grpcapi implements gRPC server API
// as a mean of communication between server and client.
// It consists of 2 services: Auth provides registration and authorization
// services, along as token refresh. Auth methods are not protected with any
// means. Docs service handle Documents operation requests. All it's methods
// require token authorization.
package grpcapi

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
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
	docs               *DocsHandlers
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	tlsCredentials     credentials.TransportCredentials
}

// GetAddress returns server binding address
func (gs *GRCPServer) GetAddress() string {
	return gs.address
}

// New is a server instance constructor
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
		grpc.Creds(gs.tlsCredentials),
		grpc.ChainUnaryInterceptor(gs.unaryInterceptors...),
		grpc.ChainStreamInterceptor(gs.streamInterceptors...),
	)

	pb.RegisterAuthServer(gs.Server, gs.auth)
	pb.RegisterDocsServer(gs.Server, gs.docs)
	return gs
}

// WithAddress allows to set server bind address
func WithAddress(s string) func(gs *GRCPServer) {
	return func(gs *GRCPServer) {
		gs.address = s
	}
}

// WithAuthHandlers defines handlers for Auth service
func WithAuthHandlers(h *AuthHandlers) func(gs *GRCPServer) {
	return func(gs *GRCPServer) {
		gs.auth = h
	}
}

// WithDocsHandlers defines handlers for Docs service
func WithDocsHandlers(h *DocsHandlers) func(gs *GRCPServer) {
	return func(gs *GRCPServer) {
		gs.docs = h
	}
}

// WithUnaryInterceptors adds unary server interceptors into the interceptors chain.
// Execution order is the same as how they were added.
func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) func(server *GRCPServer) {
	return func(gs *GRCPServer) {
		// order makes sense
		gs.unaryInterceptors = append(gs.unaryInterceptors, interceptors...)
	}
}

// WithStreamInterceptors adds stream server interceptors into the interceptors chain.
// Execution order is the same as how they were added.
func WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) func(server *GRCPServer) {
	return func(gs *GRCPServer) {
		// order makes sense
		gs.streamInterceptors = append(gs.streamInterceptors, interceptors...)
	}
}

// WithTransportCredentials sets up secure server connection
func WithTransportCredentials(cred credentials.TransportCredentials) func(server *GRCPServer) {
	return func(gs *GRCPServer) {
		gs.tlsCredentials = cred
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
