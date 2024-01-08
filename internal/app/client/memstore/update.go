package memstore

import (
	"log"

	"yap-pwkeeper/internal/pkg/models"
)

func (s *Store) getSerial() int64 {
	s.mu.RLock()
	log.Println("Rlocked")
	defer func() {
		s.mu.RUnlock()
		log.Println("Runlocked")
	}()
	return s.serial
}

func (s *Store) Update() error {
	chData := make(chan interface{})
	chErr := make(chan error, 1)
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
			if data.(models.Note).Serial > serial {
				serial = data.(models.Note).Serial
			}
		case models.Credential:
			s.placeCredential(data.(models.Credential))
			if data.(models.Credential).Serial > serial {
				serial = data.(models.Credential).Serial
			}
		case models.Card:
			s.placeCard(data.(models.Card))
			if data.(models.Card).Serial > serial {
				serial = data.(models.Card).Serial
			}
		}
	}
	err := <-chErr
	if err == nil {
		s.mu.Lock()
		log.Println("locked")
		s.serial = serial
		s.mu.Unlock()
		log.Println("unlocked")
	}
	return err
}

func (s *Store) placeNote(d models.Note) {
	s.mu.Lock()
	log.Println("locked")
	defer func() {
		s.mu.Unlock()
		log.Println("unlocked")
	}()
	if d.State == models.StateDeleted {
		delete(s.notes, d.Id)
	} else {
		s.notes[d.Id] = &d
	}
	//if d.Serial > s.serial {
	//	s.serial = d.Serial
	//}
}

func (s *Store) placeCard(d models.Card) {
	s.mu.Lock()
	log.Println("locked")
	defer func() {
		s.mu.Unlock()
		log.Println("unlocked")
	}()
	if d.State == models.StateDeleted {
		delete(s.cards, d.Id)
	} else {
		s.cards[d.Id] = &d
	}
	//if d.Serial > s.serial {
	//	s.serial = d.Serial
	//}
}

func (s *Store) placeCredential(d models.Credential) {
	s.mu.Lock()
	log.Println("locked")
	defer func() {
		s.mu.Unlock()
		log.Println("unlocked")
	}()
	if d.State == models.StateDeleted {
		delete(s.credentials, d.Id)
	} else {
		s.credentials[d.Id] = &d
	}
	//if d.Serial > s.serial {
	//	s.serial = d.Serial
	//}
}
