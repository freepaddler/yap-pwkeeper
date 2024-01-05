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
		Name:     x.Name,
		Text:     x.Text,
		Metadata: toMetadata(x.Metadata),
	}, nil
}
