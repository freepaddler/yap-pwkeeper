package grpcapi

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"

	"yap-pwkeeper/internal/pkg/grpc/interceptors"
	pb "yap-pwkeeper/internal/pkg/grpc/proto"
	"yap-pwkeeper/internal/pkg/logger"
)

type GRCPServer struct {
	Server             *grpc.Server
	address            string
	authHandlers       *AuthHandlers
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
		logging.StreamServerInterceptor(interceptors.ZapLogger(logger.Log().Desugar()), logOpts...),
	)
	for _, o := range opts {
		o(gs)
	}
	//if gs.encoder != nil {
	//	encoding.RegisterCodec(gs.encoder)
	//}
	gs.Server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(gs.unaryInterceptors...),
		grpc.ChainStreamInterceptor(gs.streamInterceptors...),
	)
	pb.RegisterAuthServer(gs.Server, gs.authHandlers)
	return gs
}

func WithAddress(s string) func(gs *GRCPServer) {
	return func(gs *GRCPServer) {
		gs.address = s
	}
}

func WithAuthHandlers(h *AuthHandlers) func(gs *GRCPServer) {
	return func(gs *GRCPServer) {
		gs.authHandlers = h
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
