package grpcapi

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"yap-pwkeeper/internal/pkg/grpc/proto"
	"yap-pwkeeper/internal/pkg/logger"
)

func (w DocsHandlers) AddCard(ctx context.Context, in *proto.Card) (*proto.Empty, error) {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("add card request")
	card, err := in.ToCard()
	card.UserId, _ = logger.GetUserId(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}
	if err := w.docs.AddCard(ctx, card); err != nil {
		return nil, respErr(ctx, err)
	}
	return &proto.Empty{}, err
}

func (w DocsHandlers) DeleteCard(ctx context.Context, in *proto.Card) (*proto.Empty, error) {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("delete card request")
	card, err := in.ToCard()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}
	card.UserId, _ = logger.GetUserId(ctx)
	if err := w.docs.DeleteCard(ctx, card); err != nil {
		return nil, respErr(ctx, err)
	}
	return &proto.Empty{}, nil
}

func (w DocsHandlers) UpdateCard(ctx context.Context, in *proto.Card) (*proto.Empty, error) {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("update card request")
	card, err := in.ToCard()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}
	card.UserId, _ = logger.GetUserId(ctx)
	if err := w.docs.UpdateCard(ctx, card); err != nil {
		return nil, respErr(ctx, err)
	}
	return &proto.Empty{}, nil
}
