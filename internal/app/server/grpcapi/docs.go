package grpcapi

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

	AddFile(ctx context.Context, file models.File) error
	DeleteFile(ctx context.Context, file models.File) error
	UpdateFile(ctx context.Context, file models.File) error
	GetFile(ctx context.Context, docId string, userId string) (models.File, error)

	GetUpdatesStream(ctx context.Context, userId string, minSerial int64, chData chan interface{}, chErr chan error)
}

type DocsHandlers struct {
	pb.UnimplementedDocsServer
	docs Docs
}

func NewDocsHandlers(db Docs) *DocsHandlers {
	return &DocsHandlers{docs: db}
}

func (w DocsHandlers) GetUpdateStream(request *pb.UpdateRequest, stream pb.Docs_GetUpdateStreamServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("update stream request")
	chData := make(chan interface{})
	chErr := make(chan error, 1)
	userId, _ := logger.GetUserId(ctx)
	go w.docs.GetUpdatesStream(ctx, userId, request.GetSerial(), chData, chErr)
	for {
		data, ok := <-chData
		if !ok {
			break
		}
		response := new(pb.UpdateResponse)
		switch data.(type) {
		case models.Note:
			response = &pb.UpdateResponse{Update: &pb.UpdateResponse_Note{Note: pb.FromNote(data.(models.Note))}}
		case models.Card:
			response = &pb.UpdateResponse{Update: &pb.UpdateResponse_Card{Card: pb.FromCard(data.(models.Card))}}
		case models.Credential:
			response = &pb.UpdateResponse{Update: &pb.UpdateResponse_Credential{Credential: pb.FromCredential(data.(models.Credential))}}
		default:
			log.Warnf("invalid data type in updates stream")
			continue
		}
		if err := stream.Send(response); err != nil {
			logger.Log().WithErr(err).Debug("update stream send failed")
			cancel()
			return status.Error(codes.Internal, "update stream send failed")
		}

	}
	if err := <-chErr; err != nil {
		return status.Error(codes.Internal, "update stream failed")
	}
	logger.Log().Debug("update stream success")
	return nil
}
