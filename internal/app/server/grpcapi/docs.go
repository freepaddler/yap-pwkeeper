package grpcapi

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"yap-pwkeeper/internal/pkg/grpc/proto"
	"yap-pwkeeper/internal/pkg/logger"
	"yap-pwkeeper/internal/pkg/models"
	"yap-pwkeeper/internal/pkg/wallet"
)

type Docs interface {
	AddNote(ctx context.Context, note models.Note) error
}

type DocsHandlers struct {
	proto.UnimplementedWalletServer
	store Docs
}

func NewDocsHandlers(db Docs) *DocsHandlers {
	return &DocsHandlers{store: db}
}

func (w DocsHandlers) AddNote(ctx context.Context, in *proto.Note) (*proto.Empty, error) {
	logger.Log().Info("addnote request")
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("add note request")
	note, err := in.ToNote()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}
	err = w.store.AddNote(ctx, note)
	switch {
	case errors.Is(wallet.ErrBadRequest, err):
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	case err != nil:
		return nil, status.Error(codes.Internal, "server error")
	}
	return &proto.Empty{}, err
}
