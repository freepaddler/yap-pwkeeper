package grpcapi

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"yap-pwkeeper/internal/pkg/grpc/proto"
	"yap-pwkeeper/internal/pkg/logger"
)

// AddNote provides AddNote document service
func (w DocsHandlers) AddNote(ctx context.Context, in *proto.Note) (*proto.Empty, error) {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("add note request")
	note, err := in.ToNote()
	note.UserId, _ = logger.GetUserId(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}
	if err := w.docs.AddNote(ctx, note); err != nil {
		return nil, respErr(ctx, err)
	}
	return &proto.Empty{}, err
}

// DeleteNote provides DeleteNote document service
func (w DocsHandlers) DeleteNote(ctx context.Context, in *proto.Note) (*proto.Empty, error) {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("delete note request")
	note, err := in.ToNote()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}
	note.UserId, _ = logger.GetUserId(ctx)
	if err := w.docs.DeleteNote(ctx, note); err != nil {
		return nil, respErr(ctx, err)
	}
	return &proto.Empty{}, nil
}

// UpdateNote provides UpdateNote document service
func (w DocsHandlers) UpdateNote(ctx context.Context, in *proto.Note) (*proto.Empty, error) {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("update note request")
	note, err := in.ToNote()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}
	note.UserId, _ = logger.GetUserId(ctx)
	if err := w.docs.UpdateNote(ctx, note); err != nil {
		return nil, respErr(ctx, err)
	}
	return &proto.Empty{}, nil
}
