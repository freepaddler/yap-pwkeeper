package memstore

import "yap-pwkeeper/internal/pkg/models"

func (s *Store) getSerial() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.serial
}

func (s *Store) Update() error {
	chData := make(chan interface{})
	chErr := make(chan error)
	serial := s.getSerial()
	go s.server.GetUpdateStream(serial, chData, chErr)
	for {
		data, ok := <-chData
		if !ok {
			break
		}
		switch data.(type) {
		case models.Note:
			s.placeNote(data.(models.Note))
		case models.Credential:
			s.placeCredential(data.(models.Credential))
		case models.Card:
			s.placeCard(data.(models.Card))
		}
	}
	err := <-chErr
	return err
}

func (s *Store) placeNote(d models.Note) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if d.State == models.StateDeleted {
		delete(s.notes, d.Id)
	} else {
		s.notes[d.Id] = &d
	}
	if d.Serial > s.serial {
		s.serial = d.Serial
	}
}

func (s *Store) placeCard(d models.Card) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if d.State == models.StateDeleted {
		delete(s.cards, d.Id)
	} else {
		s.cards[d.Id] = &d
	}
	if d.Serial > s.serial {
		s.serial = d.Serial
	}
}

func (s *Store) placeCredential(d models.Credential) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if d.State == models.StateDeleted {
		delete(s.credentials, d.Id)
	} else {
		s.credentials[d.Id] = &d
	}
	if d.Serial > s.serial {
		s.serial = d.Serial
	}
}
