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

func TestController_AddCredential(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	docStore := mocks.NewMockDocStorage(mockController)

	c := New(docStore)
	someErr := fake.Error()

	tests := []struct {
		name       string
		credential models.Credential
		retErr     error
		wantErr    error
	}{
		{
			name:       "ok",
			credential: models.Credential{},
			retErr:     nil,
			wantErr:    nil,
		},
		{
			name:       "error",
			credential: models.Credential{},
			retErr:     someErr,
			wantErr:    someErr,
		},
	}
	serial.SetSource(new(serial.SimpleSerialSource))
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, fake.Struct(&tt.credential))
			docStore.EXPECT().AddCredential(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, doc models.Credential) (string, error) {
					s, _ := serial.Next(ctx)
					assert.Equal(t, s-1, doc.Serial, "serial should be from serials source")
					assert.Equal(t, models.StateActive, doc.State, "state should be active")
					assert.Equal(t, "", doc.Id, "id should be empty")
					return fake.Word(), tt.retErr
				}).Times(1)
			err := c.AddCredential(ctx, tt.credential)
			require.ErrorIs(t, err, tt.wantErr, "expect error %s, got %s", tt.wantErr, err)
		})
	}
}

func TestController_DeleteCredential(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	docStore := mocks.NewMockDocStorage(mockController)

	c := New(docStore)
	someErr := fake.Error()

	tests := []struct {
		name             string
		credential       models.Credential
		credentialSerial int64
		storedSerial     int64
		storedState      string
		findErr          error
		retErr           error
		wantErr          error
		wantCalls        int
	}{
		{
			name:             "ok",
			credential:       models.Credential{},
			credentialSerial: 128,
			storedSerial:     128,
			storedState:      models.StateActive,
			findErr:          nil,
			retErr:           nil,
			wantErr:          nil,
			wantCalls:        1,
		},
		{
			name:             "update error",
			credential:       models.Credential{},
			credentialSerial: 128,
			storedSerial:     128,
			storedState:      models.StateActive,
			findErr:          nil,
			retErr:           someErr,
			wantErr:          someErr,
			wantCalls:        1,
		},
		{
			name:             "not found",
			credential:       models.Credential{},
			credentialSerial: 128,
			storedSerial:     128,
			storedState:      models.StateActive,
			findErr:          ErrNotFound,
			retErr:           someErr,
			wantErr:          ErrNotFound,
		},
		{
			name:             "already deleted",
			credential:       models.Credential{},
			credentialSerial: 128,
			storedSerial:     128,
			storedState:      models.StateDeleted,
			findErr:          nil,
			retErr:           someErr,
			wantErr:          ErrDeleted,
		},
		{
			name:             "serial mismatch",
			credential:       models.Credential{},
			credentialSerial: 128,
			storedSerial:     128 + 1,
			storedState:      models.StateActive,
			findErr:          nil,
			retErr:           someErr,
			wantErr:          ErrChanged,
		},
	}
	serial.SetSource(new(serial.SimpleSerialSource))
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, fake.Struct(&tt.credential))
			tt.credential.Serial = tt.credentialSerial
			docStore.EXPECT().ModifyCredential(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, doc models.Credential) error {
					s, _ := serial.Next(ctx)
					assert.Equal(t, models.Credential{
						Id:     doc.Id,
						UserId: doc.UserId,
						Serial: s - 1,
						Name:   doc.Name,
						State:  models.StateDeleted,
					}, doc)
					return tt.retErr
				}).Times(tt.wantCalls)
			docStore.EXPECT().GetCredential(ctx, tt.credential.Id, tt.credential.UserId).
				Return(models.Credential{Serial: tt.storedSerial, State: tt.storedState}, tt.findErr).Times(1)
			err := c.DeleteCredential(ctx, tt.credential)
			require.ErrorIs(t, err, tt.wantErr, "expect error %s, got %s", tt.wantErr, err)
		})
	}
}

func TestController_UpdateCredential(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	docStore := mocks.NewMockDocStorage(mockController)

	c := New(docStore)
	someErr := fake.Error()

	tests := []struct {
		name             string
		credential       models.Credential
		credentialSerial int64
		storedSerial     int64
		storedState      string
		findErr          error
		retErr           error
		wantErr          error
		wantCalls        int
	}{
		{
			name:             "ok",
			credential:       models.Credential{},
			credentialSerial: 128,
			storedSerial:     128,
			storedState:      models.StateActive,
			findErr:          nil,
			retErr:           nil,
			wantErr:          nil,
			wantCalls:        1,
		},
		{
			name:             "update error",
			credential:       models.Credential{},
			credentialSerial: 128,
			storedSerial:     128,
			storedState:      models.StateActive,
			findErr:          nil,
			retErr:           someErr,
			wantErr:          someErr,
			wantCalls:        1,
		},
		{
			name:             "not found",
			credential:       models.Credential{},
			credentialSerial: 128,
			storedSerial:     128,
			storedState:      models.StateActive,
			findErr:          ErrNotFound,
			retErr:           someErr,
			wantErr:          ErrNotFound,
		},
		{
			name:             "already deleted",
			credential:       models.Credential{},
			credentialSerial: 128,
			storedSerial:     128,
			storedState:      models.StateDeleted,
			findErr:          nil,
			retErr:           someErr,
			wantErr:          ErrDeleted,
		},
		{
			name:             "serial mismatch",
			credential:       models.Credential{},
			credentialSerial: 128,
			storedSerial:     128 + 1,
			storedState:      models.StateActive,
			findErr:          nil,
			retErr:           someErr,
			wantErr:          ErrChanged,
		},
	}
	serial.SetSource(new(serial.SimpleSerialSource))
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, fake.Struct(&tt.credential))
			tt.credential.Serial = tt.credentialSerial
			docStore.EXPECT().ModifyCredential(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, doc models.Credential) error {
					s, _ := serial.Next(ctx)
					assert.Equal(t, s-1, doc.Serial, "expect correct serial")
					assert.Equal(t, models.StateActive, doc.State, "expect state active")
					return tt.retErr
				}).Times(tt.wantCalls)
			docStore.EXPECT().GetCredential(ctx, tt.credential.Id, tt.credential.UserId).
				Return(models.Credential{Serial: tt.storedSerial, State: tt.storedState}, tt.findErr).Times(1)
			err := c.UpdateCredential(ctx, tt.credential)
			require.ErrorIs(t, err, tt.wantErr, "expect error %s, got %s", tt.wantErr, err)
		})
	}
}
