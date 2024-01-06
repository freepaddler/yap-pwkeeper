package wallet

import (
	"context"
	"errors"

	"yap-pwkeeper/internal/pkg/models"
	"yap-pwkeeper/internal/pkg/namedq"
)

var (
	ErrDeleted    = errors.New("document already deleted")
	ErrBadRequest = errors.New("invalid request data")
	ErrNotFound   = errors.New("document not found")
	ErrChanged    = errors.New("document server version mismatch")
)

type DocStorage interface {
	AddNote(ctx context.Context, note models.Note) (string, error)
	GetNote(ctx context.Context, docId string, userId string) (models.Note, error)
	ModifyNote(ctx context.Context, note models.Note) error

	AddCard(ctx context.Context, card models.Card) (string, error)
	GetCard(ctx context.Context, docId string, userId string) (models.Card, error)
	ModifyCard(ctx context.Context, card models.Card) error

	AddCredential(ctx context.Context, credential models.Credential) (string, error)
	GetCredential(ctx context.Context, docId string, userId string) (models.Credential, error)
	ModifyCredential(ctx context.Context, credential models.Credential) error
}

type Controller struct {
	store DocStorage
	queue *namedq.NamedQ
}

func New(store DocStorage) *Controller {
	c := &Controller{
		store: store,
		queue: namedq.New(),
	}
	return c
}
