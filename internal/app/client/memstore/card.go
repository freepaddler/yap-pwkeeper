package memstore

import (
	"sort"

	"yap-pwkeeper/internal/pkg/models"
)

// GetCard returns Card from store
func (s *Store) GetCard(id string) *models.Card {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cards[id]
}

// GetCardsList returns Cards array from store sorted by Name
func (s *Store) GetCardsList() []*models.Card {
	list := make([]*models.Card, 0, len(s.cards))
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.cards {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
	return list
}

// AddCard saves new Card to server
func (s *Store) AddCard(d models.Card) error {
	return s.checkAuthErr(s.server.AddCard(d))
}

// UpdateCard updates Card on server
func (s *Store) UpdateCard(d models.Card) error {
	return s.checkAuthErr(s.server.UpdateCard(d))
}

// DeleteCard deletes Card on server
func (s *Store) DeleteCard(d models.Card) error {
	return s.checkAuthErr(s.server.DeleteCard(d))
}
