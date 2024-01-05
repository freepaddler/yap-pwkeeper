package wallet

import (
	"context"
	"errors"
	"time"

	"yap-pwkeeper/internal/pkg/logger"
	"yap-pwkeeper/internal/pkg/models"
)

var (
	ErrBadRequest = errors.New("invalid request data")
	ErrNotFound   = errors.New("document not found")
)

type DocStorage interface {
	AddNote(ctx context.Context, note models.NoteDocument) error
}

type Controller struct {
	store DocStorage
}

func New(store DocStorage) *Controller {
	c := &Controller{store: store}
	return c
}

func (c *Controller) AddNote(ctx context.Context, note models.Note) error {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	userId, ok := logger.GetUserId(ctx)
	if !ok {
		userId = "qqqwww111"
		//return ErrBadRequest
	}
	log.Debug("add note request")
	nd := models.NoteDocument{
		Note: note,
		Document: models.Document{
			UserId: userId,
			Entity: models.Entity{
				CreatedAt:  time.Now(),
				ModifiedAt: time.Now(),
				State:      models.StateActive,
			},
		},
	}
	err := c.store.AddNote(ctx, nd)
	if err != nil {
		logger.Log().Warnf("add note failed: %s", err.Error())
	}
	logger.Log().Info("add note success: %s")
	return err
}
