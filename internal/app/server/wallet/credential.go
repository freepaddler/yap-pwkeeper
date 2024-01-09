package wallet

import (
	"context"

	"yap-pwkeeper/internal/app/server/serial"
	"yap-pwkeeper/internal/pkg/logger"
	"yap-pwkeeper/internal/pkg/models"
)

func (c *Controller) AddCredential(ctx context.Context, credential models.Credential) error {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("add credential request")
	qLock := c.queue.Reserve(credential.UserId)
	defer qLock.Release()
	s, err := serial.Next(ctx)
	if err != nil {
		return err
	}
	credential.Serial = s
	credential.State = models.StateActive
	credential.Id = ""
	oid, err := c.store.AddCredential(ctx, credential)
	if err != nil {
		logger.Log().Warnf("add credential failed: %s", err.Error())
	} else {
		log.With("documentId", oid).Info("credential added")
	}
	return err
}

func (c *Controller) DeleteCredential(ctx context.Context, credential models.Credential) error {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx).With("documentId", credential.Id)
	log.Debug("delete credential request")

	qLock := c.queue.Reserve(credential.UserId)
	defer qLock.Release()

	if err := c.validateCredentialUpdate(ctx, credential); err != nil {
		return err
	}

	s, err := serial.Next(ctx)
	if err != nil {
		return err
	}
	credential.Serial = s

	deleted := models.Credential{
		Id:     credential.Id,
		UserId: credential.UserId,
		Name:   credential.Name,
		Serial: s,
		State:  models.StateDeleted,
	}
	err = c.store.ModifyCredential(ctx, deleted)
	if err != nil {
		logger.Log().Warnf("credential delete failed: %s", err.Error())
	} else {
		logger.Log().Info("credential deleted")
	}
	return err
}

func (c *Controller) UpdateCredential(ctx context.Context, credential models.Credential) error {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx).With("documentId", credential.Id)
	log.Debug("update credential request")

	qLock := c.queue.Reserve(credential.UserId)
	defer qLock.Release()

	if err := c.validateCredentialUpdate(ctx, credential); err != nil {
		return err
	}

	s, err := serial.Next(ctx)
	if err != nil {
		return err
	}
	credential.Serial = s
	credential.State = models.StateActive

	err = c.store.ModifyCredential(ctx, credential)
	if err != nil {
		logger.Log().Warnf("credential update failed: %s", err.Error())
	} else {
		logger.Log().Info("credential updated")
	}
	return err
}

func (c *Controller) validateCredentialUpdate(ctx context.Context, credential models.Credential) error {
	// get stored credential
	stored, err := c.store.GetCredential(ctx, credential.Id, credential.UserId)
	if err != nil {
		return err
	}
	if stored.State == models.StateDeleted {
		return ErrDeleted
	}
	if stored.Serial > credential.Serial {
		return ErrChanged
	}
	return nil
}
