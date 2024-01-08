package memstore

import (
	"errors"
	"sync"

	"yap-pwkeeper/internal/pkg/models"
)

type DocServer interface {
	GetUpdateStream(serial int64, chData chan interface{}, chErr chan error)
}

var (
	ErrAuthFail = errors.New("authorization failed, login required")
)

type Store struct {
	notes       map[string]*models.Note
	cards       map[string]*models.Card
	credentials map[string]*models.Credential
	serial      int64
	mu          sync.RWMutex
	server      DocServer
}

func New(server DocServer) *Store {
	return &Store{
		serial:      -1,
		notes:       make(map[string]*models.Note),
		cards:       make(map[string]*models.Card),
		credentials: make(map[string]*models.Credential),
		server:      server,
	}
}
