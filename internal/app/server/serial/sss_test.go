package serial

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimpleSerialSource_GetSerials(t *testing.T) {
	sss := new(SimpleSerialSource)
	batch := 15
	prev, _ := sss.GetSerials(context.Background(), batch)
	next, _ := sss.GetSerials(context.Background(), batch)
	require.Equal(t, int64(batch), next-prev, "expect serial diff to match batch")
}
