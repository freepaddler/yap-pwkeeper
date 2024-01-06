package proto

import (
	"errors"

	"yap-pwkeeper/internal/pkg/models"
)

var (
	ErrBadRequest = errors.New("invalid request")
)

func toMetadata(x []*Meta) []models.Meta {
	metadata := make([]models.Meta, len(x))
	for i, v := range x {
		metadata[i] = models.Meta{
			Key:   v.Key,
			Value: v.Value,
		}
	}
	return metadata
}

func (x *Note) ToNote() (models.Note, error) {
	if x.Name == "" {
		return models.Note{}, ErrBadRequest
	}
	return models.Note{
		Id:       x.Id,
		Serial:   x.Serial,
		State:    x.State,
		Name:     x.Name,
		Text:     x.Text,
		Metadata: toMetadata(x.Metadata),
	}, nil
}

func (x *Credential) ToCredential() (models.Credential, error) {
	if x.Name == "" {
		return models.Credential{}, ErrBadRequest
	}
	return models.Credential{
		Id:       x.Id,
		Serial:   x.Serial,
		State:    x.State,
		Name:     x.Name,
		Login:    x.Login,
		Password: x.Password,
		Metadata: toMetadata(x.Metadata),
	}, nil
}

func (x *Card) ToCard() (models.Card, error) {
	if x.Name == "" {
		return models.Card{}, ErrBadRequest
	}
	return models.Card{
		Id:         x.Id,
		Serial:     x.Serial,
		State:      x.State,
		Name:       x.Name,
		Cardholder: x.Cardholder,
		Number:     x.Number,
		Expires:    x.Expires,
		Pin:        x.Pin,
		Code:       x.Code,
		Metadata:   toMetadata(x.Metadata),
	}, nil
}
