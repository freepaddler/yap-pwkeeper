package wallet

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"

	"yap-pwkeeper/internal/app/server/serial"
	"yap-pwkeeper/internal/pkg/logger"
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
	GetNotesStream(ctx context.Context, userId string, minSerial, maxSerial int64, chData chan interface{}) error

	AddCard(ctx context.Context, card models.Card) (string, error)
	GetCard(ctx context.Context, docId string, userId string) (models.Card, error)
	ModifyCard(ctx context.Context, card models.Card) error
	GetCardsStream(ctx context.Context, userId string, minSerial, maxSerial int64, chData chan interface{}) error

	AddCredential(ctx context.Context, credential models.Credential) (string, error)
	GetCredential(ctx context.Context, docId string, userId string) (models.Credential, error)
	ModifyCredential(ctx context.Context, credential models.Credential) error
	GetCredentialsStream(ctx context.Context, userId string, minSerial, maxSerial int64, chData chan interface{}) error
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

func (c *Controller) GetUpdatesStream(ctx context.Context, userId string, minSerial int64, chData chan interface{}, chErr chan error) {
	defer func() {
		close(chData)
		close(chErr)
	}()
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("updates stream request")
	maxSerial, err := serial.Next(ctx)
	if err != nil {
		chErr <- fmt.Errorf("unable to get new serial: %w", err)
		return
	}
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return c.store.GetNotesStream(gCtx, userId, minSerial, maxSerial, chData)
	})
	g.Go(func() error {
		return c.store.GetCardsStream(gCtx, userId, minSerial, maxSerial, chData)
	})
	g.Go(func() error {
		return c.store.GetCredentialsStream(gCtx, userId, minSerial, maxSerial, chData)
	})

	err = g.Wait()
	if err != nil {
		log.WithErr(err).Error("updates stream failed")
		chErr <- err
	}
}
