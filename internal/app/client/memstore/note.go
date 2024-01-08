package memstore

import (
	"log"
	"sort"

	"yap-pwkeeper/internal/pkg/models"
)

func (s *Store) GetNote(id string) *models.Note {
	s.mu.RLock()
	log.Println("Rlocked")
	defer func() {
		s.mu.RUnlock()
		log.Println("Runlocked")
	}()
	return s.notes[id]
}

func (s *Store) GetNotesList() []*models.Note {
	list := make([]*models.Note, 0, len(s.notes))
	s.mu.RLock()
	log.Println("Rlocked")
	defer func() {
		s.mu.RUnlock()
		log.Println("Runlocked")
	}()
	for _, v := range s.notes {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
	return list
}

func (s *Store) AddNote(note models.Note) error {
	return s.server.AddNote(note)
}

func (s *Store) UpdateNote(note models.Note) error {
	return s.server.UpdateNote(note)
}

func (s *Store) DeleteNote(note models.Note) error {
	return s.server.DeleteNote(note)
}
