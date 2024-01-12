// Package namedq helps to keep update requests of each user in one consecutive queue.
// This guarantees, that all changes of one user will not be concurrent and all updates
// will not miss any change. Also, it allows make subsequent  requests to database
// without need in transaction.
package namedq

import (
	"sync"
)

type NamedQ struct {
	mu           sync.Mutex
	reservations map[string]chan struct{}
}

// New is a constructor
func New() *NamedQ {
	return &NamedQ{
		reservations: make(map[string]chan struct{}),
	}
}

// Reservation implements Release method
type Reservation struct {
	id    string
	ch    chan struct{}
	queue *NamedQ
}

// Reserve fixes the place of requester in queue. One queue per each userId.
// As far as one user can make too many requests from different devices at
// the same time, this solution is appropriate.
// Reserved should be called before any change action.
func (nq *NamedQ) Reserve(id string) Reservation {
	ch := make(chan struct{})
	nq.mu.Lock()
	wait, ok := nq.reservations[id]
	nq.reservations[id] = ch
	nq.mu.Unlock()
	if ok {
		<-wait
	}
	return Reservation{
		id:    id,
		ch:    ch,
		queue: nq,
	}
}

// Release releases callee from the queue. Should be called after change action is complete.
func (r Reservation) Release() {
	close(r.ch)
	r.queue.mu.Lock()
	if r.ch == r.queue.reservations[r.id] {
		delete(r.queue.reservations, r.id)
	}
	defer r.queue.mu.Unlock()
}
