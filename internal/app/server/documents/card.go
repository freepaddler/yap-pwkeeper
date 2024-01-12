package documents

import (
	"context"

	"yap-pwkeeper/internal/app/server/serial"
	"yap-pwkeeper/internal/pkg/logger"
	"yap-pwkeeper/internal/pkg/models"
)

// AddCard stores new Card in DataStorage
func (c *Controller) AddCard(ctx context.Context, card models.Card) error {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("add card request")
	qLock := c.queue.Reserve(card.UserId)
	defer qLock.Release()
	s, err := serial.Next(ctx)
	if err != nil {
		return err
	}
	card.Serial = s
	card.State = models.StateActive
	card.Id = ""
	oid, err := c.store.AddCard(ctx, card)
	if err != nil {
		logger.Log().Warnf("add card failed: %s", err.Error())
	} else {
		log.With("documentId", oid).Info("card added")
	}
	return err
}

// DeleteCard removes  Card from DataStorage. Actually only document payload is deleted,
// but document id stays in DataStorage with Deleted flag
func (c *Controller) DeleteCard(ctx context.Context, card models.Card) error {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx).With("documentId", card.Id)
	log.Debug("delete card request")

	qLock := c.queue.Reserve(card.UserId)
	defer qLock.Release()

	if err := c.validateCardUpdate(ctx, card); err != nil {
		return err
	}

	s, err := serial.Next(ctx)
	if err != nil {
		return err
	}
	card.Serial = s

	deleted := models.Card{
		Id:     card.Id,
		UserId: card.UserId,
		Name:   card.Name,
		Serial: s,
		State:  models.StateDeleted,
	}
	err = c.store.ModifyCard(ctx, deleted)
	if err != nil {
		logger.Log().Warnf("card delete failed: %s", err.Error())
	} else {
		logger.Log().Info("card deleted")
	}
	return err
}

// UpdateCard modifies the whole Card, leaving id intact.
func (c *Controller) UpdateCard(ctx context.Context, card models.Card) error {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx).With("documentId", card.Id)
	log.Debug("update card request")

	qLock := c.queue.Reserve(card.UserId)
	defer qLock.Release()

	if err := c.validateCardUpdate(ctx, card); err != nil {
		return err
	}

	s, err := serial.Next(ctx)
	if err != nil {
		return err
	}
	card.Serial = s
	card.State = models.StateActive

	err = c.store.ModifyCard(ctx, card)
	if err != nil {
		logger.Log().Warnf("card update failed: %s", err.Error())
	} else {
		logger.Log().Info("card updated")
	}
	return err
}

func (c *Controller) validateCardUpdate(ctx context.Context, card models.Card) error {
	// get stored card
	stored, err := c.store.GetCard(ctx, card.Id, card.UserId)
	if err != nil {
		return err
	}
	if stored.State == models.StateDeleted {
		return ErrDeleted
	}
	if stored.Serial > card.Serial {
		return ErrChanged
	}
	return nil
}
