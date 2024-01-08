package memstore

import (
	"sync"

	"yap-pwkeeper/internal/pkg/models"
)

type DocServer interface {
	GetUpdateStream(serial int64, chData chan interface{}, chErr chan error)
	AddNote(note models.Note) error
	UpdateNote(d models.Note) error
	DeleteNote(d models.Note) error
}

type Store struct {
	notes       map[string]*models.Note
	cards       map[string]*models.Card
	credentials map[string]*models.Credential
	serial      int64
	mu          sync.RWMutex
	server      DocServer
}

func New(server DocServer) *Store {
	store := &Store{
		server: server,
	}
	store.Clear()
	return store
}

func (s *Store) Clear() {
	s.notes = make(map[string]*models.Note)
	s.cards = make(map[string]*models.Card)
	s.credentials = make(map[string]*models.Credential)
	s.serial = -1
}
