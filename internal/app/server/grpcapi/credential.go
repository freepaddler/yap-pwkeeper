package grpcapi

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"yap-pwkeeper/internal/pkg/grpc/proto"
	"yap-pwkeeper/internal/pkg/logger"
)

// AddCredential provides AddCredential document service
func (w DocsHandlers) AddCredential(ctx context.Context, in *proto.Credential) (*proto.Empty, error) {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("add credential request")
	credential, err := in.ToCredential()
	credential.UserId, _ = logger.GetUserId(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}
	if err := w.docs.AddCredential(ctx, credential); err != nil {
		return nil, respErr(ctx, err)
	}
	return &proto.Empty{}, err
}

// DeleteCredential provides DeleteCredential document service
func (w DocsHandlers) DeleteCredential(ctx context.Context, in *proto.Credential) (*proto.Empty, error) {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("delete credential request")
	credential, err := in.ToCredential()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}
	credential.UserId, _ = logger.GetUserId(ctx)
	if err := w.docs.DeleteCredential(ctx, credential); err != nil {
		return nil, respErr(ctx, err)
	}
	return &proto.Empty{}, nil
}

// UpdateCredential provides UpdateCredential document service
func (w DocsHandlers) UpdateCredential(ctx context.Context, in *proto.Credential) (*proto.Empty, error) {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("update credential request")
	credential, err := in.ToCredential()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}
	credential.UserId, _ = logger.GetUserId(ctx)
	if err := w.docs.UpdateCredential(ctx, credential); err != nil {
		return nil, respErr(ctx, err)
	}
	return &proto.Empty{}, nil
}
