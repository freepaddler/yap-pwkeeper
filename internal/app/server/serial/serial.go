package serial

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"yap-pwkeeper/internal/pkg/logger"
)

type Source interface {
	GetSerials(ctx context.Context, n int) (int64, error)
}

var (
	ErrNoSource = errors.New("unable to get serials from source")
	serial      int64 // serial to return
	maxSerial   int64 // highest serial to return
	batchSize   = 1   // number of serials to receive from store
	source      Source
	mu          sync.Mutex
)

func SetSource(src Source) {
	source = src
}

func SetBatchSize(size int) {
	batchSize = size
}

func Next(ctx context.Context) (int64, error) {
	mu.Lock()
	defer mu.Unlock()
	serial++
	if serial >= maxSerial {
		logger.Log().Debug("request new serials batch")
		if err := getNew(ctx); err != nil {
			return serial, err
		}
	}
	return serial, nil
}

func getNew(ctx context.Context) error {
	var err error
	if serial, err = source.GetSerials(ctx, batchSize); err != nil {
		return fmt.Errorf("%w: %w", ErrNoSource, err)
	}
	maxSerial = serial + int64(batchSize)
	return nil
}
