package proto

import (
	"errors"
	"fmt"

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

func fromMetadata(x []models.Meta) []*Meta {
	metadata := make([]*Meta, len(x))
	for i, v := range x {
		metadata[i] = &Meta{
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

func FromNote(x models.Note) *Note {
	return &Note{
		Id:       x.Id,
		Serial:   x.Serial,
		State:    x.State,
		Name:     x.Name,
		Text:     x.Text,
		Metadata: fromMetadata(x.Metadata),
	}
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

func FromCredential(x models.Credential) *Credential {
	return &Credential{
		Id:       x.Id,
		Serial:   x.Serial,
		State:    x.State,
		Name:     x.Name,
		Login:    x.Login,
		Password: x.Password,
		Metadata: fromMetadata(x.Metadata),
	}
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

func FromCard(x models.Card) *Card {
	return &Card{
		Id:         x.Id,
		Serial:     x.Serial,
		State:      x.State,
		Name:       x.Name,
		Cardholder: x.Cardholder,
		Number:     x.Number,
		Expires:    x.Expires,
		Pin:        x.Pin,
		Code:       x.Code,
		Metadata:   fromMetadata(x.Metadata),
	}
}

func (x *File) ToFile() (models.File, error) {
	if x.Name == "" {
		return models.File{}, fmt.Errorf("%w: Name is empty", ErrBadRequest)
	}
	if x.Filename == "" {
		return models.File{}, fmt.Errorf("%w: Fileame is empty", ErrBadRequest)
	}
	file := models.File{
		Id:       x.Id,
		Serial:   x.Serial,
		Name:     x.Name,
		Filename: x.Filename,
		Size:     x.Size,
		Data:     make([]byte, 0),
		Metadata: toMetadata(x.Metadata),
	}
	return file, nil
}

func FromFile(x models.File) *File {
	return &File{
		Id:       x.Id,
		Serial:   x.Serial,
		State:    x.State,
		Name:     x.Name,
		Filename: x.Filename,
		Size:     x.Size,
		Metadata: fromMetadata(x.Metadata),
	}
}
