package namedq

import (
	"sync"
)

type NamedQ struct {
	mu           sync.Mutex
	reservations map[string]chan struct{}
}

func New() *NamedQ {
	return &NamedQ{
		reservations: make(map[string]chan struct{}),
	}
}

type Reservation struct {
	id    string
	ch    chan struct{}
	queue *NamedQ
}

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

func (r Reservation) Release() {
	close(r.ch)
	r.queue.mu.Lock()
	if r.ch == r.queue.reservations[r.id] {
		delete(r.queue.reservations, r.id)
	}
	defer r.queue.mu.Unlock()
}
