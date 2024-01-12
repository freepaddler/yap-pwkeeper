package documents

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yap-pwkeeper/internal/app/server/serial"
	"yap-pwkeeper/internal/pkg/models"
	"yap-pwkeeper/mocks"
)

func TestController_AddNote(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	docStore := mocks.NewMockDocStorage(mockController)

	c := New(docStore)
	someErr := fake.Error()

	tests := []struct {
		name    string
		note    models.Note
		retErr  error
		wantErr error
	}{
		{
			name:    "ok",
			note:    models.Note{},
			retErr:  nil,
			wantErr: nil,
		},
		{
			name:    "error",
			note:    models.Note{},
			retErr:  someErr,
			wantErr: someErr,
		},
	}
	serial.SetSource(new(serial.SimpleSerialSource))
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, fake.Struct(&tt.note))
			docStore.EXPECT().AddNote(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, doc models.Note) (string, error) {
					s, _ := serial.Next(ctx)
					assert.Equal(t, s-1, doc.Serial, "serial should be from serials source")
					assert.Equal(t, models.StateActive, doc.State, "state should be active")
					assert.Equal(t, "", doc.Id, "id should be empty")
					return fake.Word(), tt.retErr
				}).Times(1)
			err := c.AddNote(ctx, tt.note)
			require.ErrorIs(t, err, tt.wantErr, "expect error %s, got %s", tt.wantErr, err)
		})
	}
}

func TestController_DeleteNote(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	docStore := mocks.NewMockDocStorage(mockController)

	c := New(docStore)
	someErr := fake.Error()

	tests := []struct {
		name         string
		note         models.Note
		noteSerial   int64
		storedSerial int64
		storedState  string
		findErr      error
		retErr       error
		wantErr      error
		wantCalls    int
	}{
		{
			name:         "ok",
			note:         models.Note{},
			noteSerial:   128,
			storedSerial: 128,
			storedState:  models.StateActive,
			findErr:      nil,
			retErr:       nil,
			wantErr:      nil,
			wantCalls:    1,
		},
		{
			name:         "update error",
			note:         models.Note{},
			noteSerial:   128,
			storedSerial: 128,
			storedState:  models.StateActive,
			findErr:      nil,
			retErr:       someErr,
			wantErr:      someErr,
			wantCalls:    1,
		},
		{
			name:         "not found",
			note:         models.Note{},
			noteSerial:   128,
			storedSerial: 128,
			storedState:  models.StateActive,
			findErr:      ErrNotFound,
			retErr:       someErr,
			wantErr:      ErrNotFound,
		},
		{
			name:         "already deleted",
			note:         models.Note{},
			noteSerial:   128,
			storedSerial: 128,
			storedState:  models.StateDeleted,
			findErr:      nil,
			retErr:       someErr,
			wantErr:      ErrDeleted,
		},
		{
			name:         "serial mismatch",
			note:         models.Note{},
			noteSerial:   128,
			storedSerial: 128 + 1,
			storedState:  models.StateActive,
			findErr:      nil,
			retErr:       someErr,
			wantErr:      ErrChanged,
		},
	}
	serial.SetSource(new(serial.SimpleSerialSource))
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, fake.Struct(&tt.note))
			tt.note.Serial = tt.noteSerial
			docStore.EXPECT().ModifyNote(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, doc models.Note) error {
					s, _ := serial.Next(ctx)
					assert.Equal(t, models.Note{
						Id:     doc.Id,
						UserId: doc.UserId,
						Serial: s - 1,
						Name:   doc.Name,
						State:  models.StateDeleted,
					}, doc)
					return tt.retErr
				}).Times(tt.wantCalls)
			docStore.EXPECT().GetNote(ctx, tt.note.Id, tt.note.UserId).
				Return(models.Note{Serial: tt.storedSerial, State: tt.storedState}, tt.findErr).Times(1)
			err := c.DeleteNote(ctx, tt.note)
			require.ErrorIs(t, err, tt.wantErr, "expect error %s, got %s", tt.wantErr, err)
		})
	}
}

func TestController_UpdateNote(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	docStore := mocks.NewMockDocStorage(mockController)

	c := New(docStore)
	someErr := fake.Error()

	tests := []struct {
		name         string
		note         models.Note
		noteSerial   int64
		storedSerial int64
		storedState  string
		findErr      error
		retErr       error
		wantErr      error
		wantCalls    int
	}{
		{
			name:         "ok",
			note:         models.Note{},
			noteSerial:   128,
			storedSerial: 128,
			storedState:  models.StateActive,
			findErr:      nil,
			retErr:       nil,
			wantErr:      nil,
			wantCalls:    1,
		},
		{
			name:         "update error",
			note:         models.Note{},
			noteSerial:   128,
			storedSerial: 128,
			storedState:  models.StateActive,
			findErr:      nil,
			retErr:       someErr,
			wantErr:      someErr,
			wantCalls:    1,
		},
		{
			name:         "not found",
			note:         models.Note{},
			noteSerial:   128,
			storedSerial: 128,
			storedState:  models.StateActive,
			findErr:      ErrNotFound,
			retErr:       someErr,
			wantErr:      ErrNotFound,
		},
		{
			name:         "already deleted",
			note:         models.Note{},
			noteSerial:   128,
			storedSerial: 128,
			storedState:  models.StateDeleted,
			findErr:      nil,
			retErr:       someErr,
			wantErr:      ErrDeleted,
		},
		{
			name:         "serial mismatch",
			note:         models.Note{},
			noteSerial:   128,
			storedSerial: 128 + 1,
			storedState:  models.StateActive,
			findErr:      nil,
			retErr:       someErr,
			wantErr:      ErrChanged,
		},
	}
	serial.SetSource(new(serial.SimpleSerialSource))
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, fake.Struct(&tt.note))
			tt.note.Serial = tt.noteSerial
			docStore.EXPECT().ModifyNote(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, doc models.Note) error {
					s, _ := serial.Next(ctx)
					assert.Equal(t, s-1, doc.Serial, "expect correct serial")
					assert.Equal(t, models.StateActive, doc.State, "expect state active")
					return tt.retErr
				}).Times(tt.wantCalls)
			docStore.EXPECT().GetNote(ctx, tt.note.Id, tt.note.UserId).
				Return(models.Note{Serial: tt.storedSerial, State: tt.storedState}, tt.findErr).Times(1)
			err := c.UpdateNote(ctx, tt.note)
			require.ErrorIs(t, err, tt.wantErr, "expect error %s, got %s", tt.wantErr, err)
		})
	}
}
