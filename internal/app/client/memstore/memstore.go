// Package memstore is a local documents cache. Cache is fully verified by server.
// No local changes are allowed. All changes are sent directly to server and are applied
// locally on receiving as updates from server side.
// Update method should be used to get server updates.
// Documents in local storage are available until server authorisation exists.
// Any unauthorised request leads to empty local storage and requires new authentication.
package memstore

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"golang.org/x/sync/singleflight"

	"yap-pwkeeper/internal/app/client/grpccli"
	"yap-pwkeeper/internal/pkg/models"
)

// DocServer is interface that implements server client
type DocServer interface {
	Register(login, password string) error
	Login(login, password string) error

	GetUpdateStream(serial int64, chData chan interface{}, chErr chan error)

	AddNote(note models.Note) error
	UpdateNote(d models.Note) error
	DeleteNote(d models.Note) error

	AddCard(note models.Card) error
	UpdateCard(d models.Card) error
	DeleteCard(d models.Card) error

	AddCredential(note models.Credential) error
	UpdateCredential(d models.Credential) error
	DeleteCredential(d models.Credential) error

	AddFile(d models.File, r io.Reader) error
	UpdateFileInfo(d models.File) error
	UpdateFile(d models.File, r io.Reader) error
	DeleteFile(d models.File) error
}

var (
	// ErrAuthFailed notifies that server authorisation is lost and no data may be accessed
	ErrAuthFailed = errors.New("authorization failed, You need to login again")
)

type Store struct {
	notes       map[string]*models.Note
	cards       map[string]*models.Card
	credentials map[string]*models.Credential
	files       map[string]*models.File
	serial      int64
	mu          sync.RWMutex
	server      DocServer
	updateGroup *singleflight.Group
}

// New is a storage constructor
func New(server DocServer) *Store {
	store := &Store{
		server:      server,
		updateGroup: new(singleflight.Group),
	}
	store.bootstrap()
	return store
}

// bootstrap creates new empty storage
func (s *Store) bootstrap() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notes = make(map[string]*models.Note)
	s.cards = make(map[string]*models.Card)
	s.credentials = make(map[string]*models.Credential)
	s.files = make(map[string]*models.File)
	s.serial = -1
}

// Register registers new server user
func (s *Store) Register(login, password string) error {
	if err := s.server.Register(login, password); err != nil {
		return fmt.Errorf("registraition failed: %w", err)
	}
	return s.Login(login, password)
}

// Login creates new server authorization session
func (s *Store) Login(login, password string) error {
	if err := s.server.Login(login, password); err != nil {
		return fmt.Errorf("login failed: %w", err)
	}
	s.bootstrap()
	return nil
}

// checkAuthErr is server response error wrapper.
// If authorised session terminates it clears storage.
func (s *Store) checkAuthErr(err error) error {
	if errors.Is(grpccli.ErrAuthFail, err) {
		s.bootstrap()
		return ErrAuthFailed
	}
	return err
}
