package serial

import (
	"context"
	"sync"
)

// SimpleSerialSource generator fo test purposes
type SimpleSerialSource struct {
	serial int64
	mu     sync.Mutex
}

func (sss *SimpleSerialSource) GetSerials(_ context.Context, n int) (int64, error) {
	sss.mu.Lock()
	defer sss.mu.Unlock()
	ret := serial
	serial = serial + int64(n)
	return ret, nil
}
