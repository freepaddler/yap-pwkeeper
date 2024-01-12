package namedq

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNamedQ(t *testing.T) {
	n := 1000
	want := make([]int, 0, n)
	got := make([]int, 0, n)
	wg := sync.WaitGroup{}
	nq := New()
	for i := 0; i < 100; i++ {
		want = append(want, i)
		i := i
		wg.Add(1)
		reserve := nq.Reserve("a")
		go func(r Reservation) {
			defer func() {
				reserve.Release()
				wg.Done()
			}()
			got = append(got, i)
		}(reserve)
	}
	wg.Wait()
	require.Equal(t, want, got, "slices expected to match")
}
