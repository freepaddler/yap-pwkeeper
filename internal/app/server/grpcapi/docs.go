package grpcapi

import (
	"context"

	pb "yap-pwkeeper/internal/pkg/grpc/proto"
	"yap-pwkeeper/internal/pkg/logger"
	"yap-pwkeeper/internal/pkg/models"
)

type Docs interface {
	AddNote(ctx context.Context, note models.Note) error
	DeleteNote(ctx context.Context, note models.Note) error
	UpdateNote(ctx context.Context, note models.Note) error

	AddCard(ctx context.Context, card models.Card) error
	DeleteCard(ctx context.Context, card models.Card) error
	UpdateCard(ctx context.Context, card models.Card) error

	AddCredential(ctx context.Context, credential models.Credential) error
	DeleteCredential(ctx context.Context, credential models.Credential) error
	UpdateCredential(ctx context.Context, credential models.Credential) error
}

type DocsHandlers struct {
	pb.UnimplementedWalletServer
	docs Docs
}

func NewDocsHandlers(db Docs) *DocsHandlers {
	return &DocsHandlers{docs: db}
}

func (w DocsHandlers) GetUpdates(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("update request")
	response := &pb.UpdateResponse{
		Update: &pb.UpdateResponse_Note{Note: &pb.Note{
			Id:     "someId",
			Serial: 0,
			State:  "active",
			Name:   "somename",
			Text:   "sometext",
		}},
	}
	log.Debug("update response")
	return response, nil
}

func (w DocsHandlers) GetUpdate(request *pb.UpdateRequest, stream pb.Wallet_GetUpdateServer) error {
	ctx := stream.Context()
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("update request")
	response := &pb.UpdateResponse{
		Update: &pb.UpdateResponse_Note{Note: &pb.Note{
			Id:     "someId",
			Serial: 0,
			State:  "active",
			Name:   "somename",
			Text:   "sometext",
		}},
	}
	if err := stream.Send(response); err != nil {
		logger.Log().WithErr(err).Debug("note send failed")
	}
	response = &pb.UpdateResponse{
		Update: &pb.UpdateResponse_Credential{
			Credential: &pb.Credential{
				Id:     "credId",
				Serial: 0,
				State:  "active",
				Name:   "credName",
				Login:  "credLogin",
			},
		},
	}
	if err := stream.Send(response); err != nil {
		logger.Log().WithErr(err).Debug("cred send failed")
	}
	logger.Log().Debug("send finished")
	return nil
}
