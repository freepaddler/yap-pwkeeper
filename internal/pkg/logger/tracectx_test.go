package logger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_UserId(t *testing.T) {
	tests := []struct {
		name      string
		wantOk    bool
		wantValue string
	}{
		{
			name:      "context with value",
			wantOk:    true,
			wantValue: "some text",
		},
		{
			name:   "context without value",
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.wantOk {
				ctx = WithUserId(ctx, tt.wantValue)
			}
			got, ok := GetUserId(ctx)
			require.Equal(t, tt.wantOk, ok)
			if tt.wantOk {
				require.Equal(t, got, tt.wantValue)
			}
		})
	}
}

func Test_RequestId(t *testing.T) {
	tests := []struct {
		name      string
		wantOk    bool
		wantValue string
	}{
		{
			name:      "context with value",
			wantOk:    true,
			wantValue: "some text",
		},
		{
			name:   "context without value",
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.wantOk {
				ctx = WithRequestId(ctx, tt.wantValue)
			}
			got, ok := GetRequestId(ctx)
			require.Equal(t, tt.wantOk, ok)
			if tt.wantOk {
				require.Equal(t, got, tt.wantValue)
			}
		})
	}
}
