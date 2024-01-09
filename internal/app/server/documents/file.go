package documents

import (
	"context"

	"yap-pwkeeper/internal/app/server/serial"
	"yap-pwkeeper/internal/pkg/logger"
	"yap-pwkeeper/internal/pkg/models"
)

func (c *Controller) GetFile(ctx context.Context, docId string, userId string) (models.File, error) {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("get file request")
	return c.GetFile(ctx, docId, userId)
}

func (c *Controller) AddFile(ctx context.Context, file models.File) error {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("add file request")
	qLock := c.queue.Reserve(file.UserId)
	defer qLock.Release()
	s, err := serial.Next(ctx)
	if err != nil {
		return err
	}
	file.Serial = s
	file.State = models.StateActive
	file.Id = ""
	oid, err := c.store.AddFile(ctx, file)
	if err != nil {
		logger.Log().Warnf("add file failed: %s", err.Error())
	} else {
		log.With("documentId", oid).Info("file added")
	}
	return err
}

func (c *Controller) DeleteFile(ctx context.Context, file models.File) error {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx).With("documentId", file.Id)
	log.Debug("delete file request")

	qLock := c.queue.Reserve(file.UserId)
	defer qLock.Release()

	if _, err := c.validateFileUpdate(ctx, file); err != nil {
		return err
	}

	s, err := serial.Next(ctx)
	if err != nil {
		return err
	}
	file.Serial = s

	deleted := models.File{
		Id:     file.Id,
		UserId: file.UserId,
		Name:   file.Name,
		Serial: s,
		State:  models.StateDeleted,
	}
	err = c.store.ModifyFile(ctx, deleted)
	if err != nil {
		logger.Log().Warnf("file delete failed: %s", err.Error())
	} else {
		logger.Log().Info("file deleted")
	}
	return err
}

func (c *Controller) UpdateFile(ctx context.Context, file models.File) error {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx).With("documentId", file.Id)
	log.Debug("update file request")

	qLock := c.queue.Reserve(file.UserId)
	defer qLock.Release()

	hash, err := c.validateFileUpdate(ctx, file)
	if err != nil {
		return err
	}

	s, err := serial.Next(ctx)
	if err != nil {
		return err
	}
	file.Serial = s
	file.State = models.StateActive

	if file.Sha265 == hash {
		err = c.store.ModifyFileInfo(ctx, file)
	} else {
		err = c.store.ModifyFile(ctx, file)
	}
	if err != nil {
		logger.Log().Warnf("file update failed: %s", err.Error())
	} else {
		logger.Log().Info("file updated")
	}
	return err
}

func (c *Controller) validateFileUpdate(ctx context.Context, file models.File) (string, error) {
	hash := ""
	// get stored file
	stored, err := c.store.GetFileInfo(ctx, file.Id, file.UserId)
	if err != nil {
		return hash, err
	}
	if stored.State == models.StateDeleted {
		return hash, ErrDeleted
	}
	if stored.Serial > file.Serial {
		return hash, ErrChanged
	}
	return hash, nil
}
