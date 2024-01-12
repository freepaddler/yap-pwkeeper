package serial

import (
	"context"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"yap-pwkeeper/mocks"
)

func TestSerial(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	src := mocks.NewMockSource(mockController)
	ctx := context.Background()
	batch := 1000
	var storedSerial int64 = 0
	var initSerial int64 = storedSerial
	reqCount := 10021
	wantCalls := reqCount/batch + 1
	wantNext := initSerial + int64(reqCount)
	SetSource(src)
	SetBatchSize(batch)

	src.EXPECT().GetSerials(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, b int) (int64, error) {
		r := storedSerial
		storedSerial = r + int64(b)
		return r, nil
	}).Times(wantCalls)

	var wg sync.WaitGroup
	wg.Add(reqCount)
	for i := 0; i < reqCount; i++ {
		go func() {
			defer wg.Done()
			_, err := Next(ctx)
			require.NoError(t, err, "no error expected on next")
		}()
	}
	wg.Wait()
	got, err := Next(ctx)
	require.NoError(t, err, "no error expected on next")
	require.Equal(t, got, wantNext, "expected next serial %d, got %d", wantNext, got)
	//fmt.Println(got, wantNext)

}
