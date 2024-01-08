package memstore

import (
	"log"
	"sort"

	"yap-pwkeeper/internal/pkg/models"
)

func (s *Store) GetCard(id string) *models.Card {
	s.mu.RLock()
	log.Println("Rlocked")
	defer func() {
		s.mu.RUnlock()
		log.Println("Runlocked")
	}()
	return s.cards[id]
}

func (s *Store) GetCardsList() []*models.Card {
	list := make([]*models.Card, 0, len(s.cards))
	s.mu.RLock()
	log.Println("Rlocked")
	defer func() {
		s.mu.RUnlock()
		log.Println("Runlocked")
	}()
	for _, v := range s.cards {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
	return list
}
