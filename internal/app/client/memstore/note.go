package memstore

import (
	"sort"

	"yap-pwkeeper/internal/pkg/models"
)

// GetNote returns Note from store
func (s *Store) GetNote(id string) *models.Note {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.notes[id]
}

// GetNotesList returns Notes array from store sorted by Name
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

// AddNote saves new Note to server
func (s *Store) AddNote(d models.Note) error {
	return s.checkAuthErr(s.server.AddNote(d))
}

// UpdateNote updates Note on server
func (s *Store) UpdateNote(d models.Note) error {
	return s.checkAuthErr(s.server.AddNote(d))
}

// DeleteNote deletes Note on server
func (s *Store) DeleteNote(d models.Note) error {
	return s.checkAuthErr(s.server.AddNote(d))
}
