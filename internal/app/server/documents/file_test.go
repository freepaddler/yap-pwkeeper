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

func TestController_AddFile(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	docStore := mocks.NewMockDocStorage(mockController)

	c := New(docStore)
	someErr := fake.Error()

	tests := []struct {
		name    string
		file    models.File
		retErr  error
		wantErr error
	}{
		{
			name:    "ok",
			file:    models.File{},
			retErr:  nil,
			wantErr: nil,
		},
		{
			name:    "error",
			file:    models.File{},
			retErr:  someErr,
			wantErr: someErr,
		},
	}
	serial.SetSource(new(serial.SimpleSerialSource))
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, fake.Struct(&tt.file))
			docStore.EXPECT().AddFile(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, doc models.File) (string, error) {
					s, _ := serial.Next(ctx)
					assert.Equal(t, s-1, doc.Serial, "serial should be from serials source")
					assert.Equal(t, models.StateActive, doc.State, "state should be active")
					assert.Equal(t, "", doc.Id, "id should be empty")
					return fake.Word(), tt.retErr
				}).Times(1)
			err := c.AddFile(ctx, tt.file)
			require.ErrorIs(t, err, tt.wantErr, "expect error %s, got %s", tt.wantErr, err)
		})
	}
}

func TestController_DeleteFile(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	docStore := mocks.NewMockDocStorage(mockController)

	c := New(docStore)
	someErr := fake.Error()

	tests := []struct {
		name         string
		file         models.File
		fileSerial   int64
		storedSerial int64
		storedState  string
		findErr      error
		retErr       error
		wantErr      error
		wantCalls    int
	}{
		{
			name:         "ok",
			file:         models.File{},
			fileSerial:   128,
			storedSerial: 128,
			storedState:  models.StateActive,
			findErr:      nil,
			retErr:       nil,
			wantErr:      nil,
			wantCalls:    1,
		},
		{
			name:         "update error",
			file:         models.File{},
			fileSerial:   128,
			storedSerial: 128,
			storedState:  models.StateActive,
			findErr:      nil,
			retErr:       someErr,
			wantErr:      someErr,
			wantCalls:    1,
		},
		{
			name:         "not found",
			file:         models.File{},
			fileSerial:   128,
			storedSerial: 128,
			storedState:  models.StateActive,
			findErr:      ErrNotFound,
			retErr:       someErr,
			wantErr:      ErrNotFound,
		},
		{
			name:         "already deleted",
			file:         models.File{},
			fileSerial:   128,
			storedSerial: 128,
			storedState:  models.StateDeleted,
			findErr:      nil,
			retErr:       someErr,
			wantErr:      ErrDeleted,
		},
		{
			name:         "serial mismatch",
			file:         models.File{},
			fileSerial:   128,
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
			require.NoError(t, fake.Struct(&tt.file))
			tt.file.Serial = tt.fileSerial
			docStore.EXPECT().ModifyFile(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, doc models.File) error {
					s, _ := serial.Next(ctx)
					assert.Equal(t, models.File{
						Id:       doc.Id,
						UserId:   doc.UserId,
						Serial:   s - 1,
						Name:     doc.Name,
						Filename: doc.Filename,
						State:    models.StateDeleted,
					}, doc)
					return tt.retErr
				}).Times(tt.wantCalls)
			docStore.EXPECT().GetFileInfo(ctx, tt.file.Id, tt.file.UserId).
				Return(models.File{Serial: tt.storedSerial, State: tt.storedState}, tt.findErr).Times(1)
			err := c.DeleteFile(ctx, tt.file)
			require.ErrorIs(t, err, tt.wantErr, "expect error %s, got %s", tt.wantErr, err)
		})
	}
}

func TestController_GetFile(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	docStore := mocks.NewMockDocStorage(mockController)

	c := New(docStore)
	someErr := fake.Error()

	tests := []struct {
		name    string
		docId   string
		userId  string
		retErr  error
		wantErr error
	}{
		{
			name:    "ok",
			docId:   fake.Word(),
			userId:  fake.Word(),
			retErr:  nil,
			wantErr: nil,
		},
		{
			name:    "failed",
			docId:   fake.Word(),
			userId:  fake.Word(),
			retErr:  someErr,
			wantErr: someErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			docStore.EXPECT().GetFile(ctx, tt.docId, tt.userId).Return(models.File{}, tt.retErr).Times(1)
			_, err := c.GetFile(ctx, tt.docId, tt.userId)
			require.ErrorIs(t, err, tt.wantErr, "expect error %s, got %s", tt.wantErr, err)
		})
	}
}

func TestController_UpdateFile(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	docStore := mocks.NewMockDocStorage(mockController)

	c := New(docStore)
	someErr := fake.Error()
	goodhash := fake.DigitN(32)
	badhash := fake.DigitN(32)

	tests := []struct {
		name                string
		file                models.File
		fileHash            string
		fileSerial          int64
		storedHash          string
		storedSerial        int64
		storedState         string
		findErr             error
		retErr              error
		wantErr             error
		wantUpdateCalls     int
		wantUpdateInfoCalls int
	}{
		{
			name:            "ok updateFile",
			file:            models.File{},
			fileSerial:      128,
			storedSerial:    128,
			storedState:     models.StateActive,
			fileHash:        goodhash,
			storedHash:      badhash,
			findErr:         nil,
			retErr:          nil,
			wantErr:         nil,
			wantUpdateCalls: 1,
		},
		{
			name:                "ok updateFileInfo",
			file:                models.File{},
			fileSerial:          128,
			storedSerial:        128,
			storedState:         models.StateActive,
			fileHash:            goodhash,
			storedHash:          goodhash,
			findErr:             nil,
			retErr:              nil,
			wantErr:             nil,
			wantUpdateInfoCalls: 1,
		},
		{
			name:            "updateFile error",
			file:            models.File{},
			fileSerial:      128,
			storedSerial:    128,
			storedState:     models.StateActive,
			fileHash:        goodhash,
			storedHash:      badhash,
			findErr:         nil,
			retErr:          someErr,
			wantErr:         someErr,
			wantUpdateCalls: 1,
		},
		{
			name:                "updateFileInfo error",
			file:                models.File{},
			fileSerial:          128,
			storedSerial:        128,
			storedState:         models.StateActive,
			fileHash:            goodhash,
			storedHash:          goodhash,
			findErr:             nil,
			retErr:              someErr,
			wantErr:             someErr,
			wantUpdateInfoCalls: 1,
		},
		{
			name:         "not found",
			file:         models.File{},
			fileSerial:   128,
			fileHash:     goodhash,
			storedSerial: 128,
			storedState:  models.StateActive,
			findErr:      ErrNotFound,
			retErr:       someErr,
			wantErr:      ErrNotFound,
		},
		{
			name:         "already deleted",
			file:         models.File{},
			fileSerial:   128,
			storedSerial: 128,
			storedState:  models.StateDeleted,
			findErr:      nil,
			retErr:       someErr,
			wantErr:      ErrDeleted,
		},
		{
			name:         "serial mismatch",
			file:         models.File{},
			fileSerial:   128,
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
			require.NoError(t, fake.Struct(&tt.file))
			tt.file.Serial = tt.fileSerial
			tt.file.Sha265 = tt.fileHash
			docStore.EXPECT().ModifyFile(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, doc models.File) error {
					s, _ := serial.Next(ctx)
					assert.Equal(t, s-1, doc.Serial, "expect correct serial")
					assert.Equal(t, models.StateActive, doc.State, "expect state active")
					return tt.retErr
				}).Times(tt.wantUpdateCalls)
			docStore.EXPECT().ModifyFileInfo(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, doc models.File) error {
					s, _ := serial.Next(ctx)
					assert.Equal(t, s-1, doc.Serial, "expect correct serial")
					assert.Equal(t, models.StateActive, doc.State, "expect state active")
					return tt.retErr
				}).Times(tt.wantUpdateInfoCalls)
			docStore.EXPECT().GetFileInfo(ctx, tt.file.Id, tt.file.UserId).
				Return(models.File{Serial: tt.storedSerial, State: tt.storedState, Sha265: tt.storedHash}, tt.findErr).Times(1)
			err := c.UpdateFile(ctx, tt.file)
			require.ErrorIs(t, err, tt.wantErr, "expect error %s, got %s", tt.wantErr, err)
		})
	}
}
