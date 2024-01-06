package grpcapi

import (
	"context"

	"yap-pwkeeper/internal/pkg/grpc/proto"
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
	proto.UnimplementedWalletServer
	docs Docs
}

func NewDocsHandlers(db Docs) *DocsHandlers {
	return &DocsHandlers{docs: db}
}
