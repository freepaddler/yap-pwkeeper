package memstore

import (
	"sort"

	"yap-pwkeeper/internal/pkg/models"
)

func (s *Store) GetNote(id string) *models.Note {
	s.mu.RLock()
	defer s.mu.RLock()
	return s.notes[id]
}

func (s *Store) GetNotesList() []*models.Note {
	list := make([]*models.Note, 0, len(s.notes))
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.notes {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
	return list
}
