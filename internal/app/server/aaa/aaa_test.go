package aaa

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"yap-pwkeeper/internal/pkg/jwtToken"
	"yap-pwkeeper/internal/pkg/models"
	"yap-pwkeeper/mocks"
)

func TestController_Register(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	userStorage := mocks.NewMockUserStorage(mockController)

	customError := errors.New("some custom error")

	tests := []struct {
		name    string
		retErr  error
		wantErr error
	}{
		{
			name:    "registered",
			retErr:  nil,
			wantErr: nil,
		},
		{
			name:    "duplicate",
			retErr:  ErrDuplicate,
			wantErr: ErrDuplicate,
		},
		{
			name:    "other error",
			retErr:  customError,
			wantErr: customError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := New(userStorage)
			ctx := context.Background()
			userStorage.EXPECT().AddUser(ctx, gomock.Any()).Return(models.User{}, tt.retErr).Times(1)
			err := controller.Register(context.Background(), models.UserCredentials{})
			assert.ErrorIs(t, err, tt.wantErr, "expected error %s, got error %s", tt.wantErr, err)
		})
	}
}

func TestController_Login(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	userStorage := mocks.NewMockUserStorage(mockController)

	customError := errors.New("some custom error")

	goodPassword := "someGoodPassword"
	goodHash, _ := bcrypt.GenerateFromPassword([]byte(goodPassword), bcrypt.DefaultCost)

	tests := []struct {
		name     string
		password string
		retHash  []byte
		retErr   error
		wantErr  error
	}{
		{
			name:     "password match",
			password: goodPassword,
			retHash:  goodHash,
			retErr:   nil,
			wantErr:  nil,
		},
		{
			name:     "password does not match",
			password: "bad password",
			retHash:  goodHash,
			retErr:   nil,
			wantErr:  ErrBadAuth,
		},
		{
			name:    "user not found",
			retErr:  ErrNotFound,
			wantErr: ErrBadAuth,
		},
		{
			name:    "other error",
			retErr:  customError,
			wantErr: customError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := New(userStorage)
			ctx := context.Background()
			userStorage.EXPECT().GetUserByLogin(ctx, gomock.Any()).
				Return(models.User{Id: "any", PasswordHash: string(tt.retHash)}, tt.retErr).Times(1)
			token, err := controller.Login(ctx, models.UserCredentials{Password: tt.password})
			if tt.wantErr == nil {
				require.NoError(t, err, "no error expected")
				require.NotEqual(t, "", token, "expected not empty token")
				return
			}
			require.ErrorIs(t, err, tt.wantErr, "expected error %s, got error %s", tt.wantErr.Error(), err.Error())
			require.Equal(t, "", token, "empty token expected")
		})
	}
}

func TestController_Refresh(t *testing.T) {
	invalid, err := jwtToken.NewToken("invalid")
	assert.NoError(t, err, "no error expected on new token request")
	jwtToken.SetKey("12345678")
	valid, err := jwtToken.NewToken("valid")
	assert.NoError(t, err, "no error expected on new token request")
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	userStorage := mocks.NewMockUserStorage(mockController)
	controller := New(userStorage)
	t.Run("valid token", func(t *testing.T) {
		token, err := controller.Refresh(context.Background(), valid)
		require.NoError(t, err, "no error expected")
		assert.NotEqual(t, "", token, "token should not be empty")
	})
	t.Run("invalid token", func(t *testing.T) {
		token, err := controller.Refresh(context.Background(), invalid)
		require.ErrorIs(t, err, ErrBadAuth, "error expected %s, got %s", err.Error(), ErrBadAuth.Error())
		assert.Equal(t, "", token, "token should be empty")
	})

}
