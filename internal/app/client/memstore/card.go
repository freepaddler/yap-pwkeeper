package memstore

import (
	"sort"

	"yap-pwkeeper/internal/pkg/models"
)

func (s *Store) GetCard(id string) *models.Card {
	s.mu.RLock()
	defer s.mu.RLock()
	return s.cards[id]
}

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
