package memstore

import (
	"log"

	"yap-pwkeeper/internal/pkg/models"
)

// Update is public update function wrapped to avoid concurrency calls
// Applies server updates to local storage.
func (s *Store) Update() error {
	_, err, _ := s.updateGroup.Do("1", func() (interface{}, error) {
		return nil, s.update()
	})
	return s.checkAuthErr(err)
}

// incSerial increases serial if new is greater than old
func incSerial(old, new int64) int64 {
	if new > old {
		return new
	}
	return old
}

// getSerial returns serial of latest update stored
func (s *Store) getSerial() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.serial
}

// update request updates from server and applies them to the store
func (s *Store) update() error {
	chData := make(chan interface{})
	chErr := make(chan error, 1)
	serial := s.getSerial()
	log.Printf("Serial before update %d", serial)
	go s.server.GetUpdateStream(serial, chData, chErr)
	for {
		data, ok := <-chData
		if !ok {
			break
		}
		switch data.(type) {
		case models.Note:
			d := data.(models.Note)
			s.placeNote(d)
			serial = incSerial(serial, d.Serial)
		case models.Credential:
			d := data.(models.Credential)
			s.placeCredential(d)
			serial = incSerial(serial, d.Serial)
		case models.Card:
			d := data.(models.Card)
			s.placeCard(d)
			serial = incSerial(serial, d.Serial)
		case models.File:
			d := data.(models.File)
			s.placeFile(d)
			serial = incSerial(serial, d.Serial)
		}
	}
	err := <-chErr
	if err == nil {
		s.mu.Lock()
		s.serial = serial
		s.mu.Unlock()
	}
	log.Printf("Serial after update %d", serial)
	return err
}

// placeNote updates or adds Note to local storage
func (s *Store) placeNote(d models.Note) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if d.State == models.StateDeleted {
		delete(s.notes, d.Id)
	} else {
		s.notes[d.Id] = &d
	}
}

// placeCard updates or adds Card to local storage
func (s *Store) placeCard(d models.Card) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if d.State == models.StateDeleted {
		delete(s.cards, d.Id)
	} else {
		s.cards[d.Id] = &d
	}
}

// placeCredential updates or adds Credentials to local storage
func (s *Store) placeCredential(d models.Credential) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if d.State == models.StateDeleted {
		delete(s.credentials, d.Id)
	} else {
		s.credentials[d.Id] = &d
	}
}

// placeFile updates or adds FileInfo to local storage
func (s *Store) placeFile(d models.File) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if d.State == models.StateDeleted {
		delete(s.files, d.Id)
	} else {
		s.files[d.Id] = &d
	}
}
