// Package serial provides requester with unique monotonic increasing serial.
// Serials are the key to the documents updates. In conjunction with namedq
// it guarantees consequent update identification and assures that no update
// action will miss any update.
package serial

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"yap-pwkeeper/internal/pkg/logger"
)

// Source is the interface to reserve next bunch of serials.
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

// SetSource allows to set source for serials
func SetSource(src Source) {
	source = src
}

// SetBatchSize configure how many serials will be reserved,
// not to query Source too frequently
func SetBatchSize(size int) {
	batchSize = size
}

// Next returns next unused serial. This is the only method, that should be
// used to get serial.
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
