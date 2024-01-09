package documents

import (
	"context"

	"yap-pwkeeper/internal/app/server/serial"
	"yap-pwkeeper/internal/pkg/logger"
	"yap-pwkeeper/internal/pkg/models"
)

func (c *Controller) AddNote(ctx context.Context, note models.Note) error {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("add note request")
	qLock := c.queue.Reserve(note.UserId)
	defer qLock.Release()
	s, err := serial.Next(ctx)
	if err != nil {
		return err
	}
	note.Serial = s
	note.State = models.StateActive
	note.Id = ""
	oid, err := c.store.AddNote(ctx, note)
	if err != nil {
		logger.Log().Warnf("add note failed: %s", err.Error())
	} else {
		log.With("documentId", oid).Info("note added")
	}
	return err
}

func (c *Controller) DeleteNote(ctx context.Context, note models.Note) error {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx).With("documentId", note.Id)
	log.Debug("delete note request")

	qLock := c.queue.Reserve(note.UserId)
	defer qLock.Release()

	if err := c.validateNoteUpdate(ctx, note); err != nil {
		return err
	}

	s, err := serial.Next(ctx)
	if err != nil {
		return err
	}
	note.Serial = s

	deleted := models.Note{
		Id:     note.Id,
		UserId: note.UserId,
		Name:   note.Name,
		Serial: s,
		State:  models.StateDeleted,
	}
	err = c.store.ModifyNote(ctx, deleted)
	if err != nil {
		logger.Log().Warnf("note delete failed: %s", err.Error())
	} else {
		logger.Log().Info("note deleted")
	}
	return err
}

func (c *Controller) UpdateNote(ctx context.Context, note models.Note) error {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx).With("documentId", note.Id)
	log.Debug("update note request")

	qLock := c.queue.Reserve(note.UserId)
	defer qLock.Release()

	if err := c.validateNoteUpdate(ctx, note); err != nil {
		return err
	}

	s, err := serial.Next(ctx)
	if err != nil {
		return err
	}
	note.Serial = s
	note.State = models.StateActive

	err = c.store.ModifyNote(ctx, note)
	if err != nil {
		logger.Log().Warnf("note update failed: %s", err.Error())
	} else {
		logger.Log().Info("note updated")
	}
	return err
}

func (c *Controller) validateNoteUpdate(ctx context.Context, note models.Note) error {
	// get stored note
	stored, err := c.store.GetNote(ctx, note.Id, note.UserId)
	if err != nil {
		return err
	}
	if stored.State == models.StateDeleted {
		return ErrDeleted
	}
	if stored.Serial > note.Serial {
		return ErrChanged
	}
	return nil
}
