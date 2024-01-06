package interceptors

import (
	"context"

	"google.golang.org/grpc"
)

type serverStreamWrapped struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *serverStreamWrapped) Context() context.Context {
	return w.ctx
}
