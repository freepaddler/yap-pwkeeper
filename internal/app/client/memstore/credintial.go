package memstore

import (
	"log"
	"sort"

	"yap-pwkeeper/internal/pkg/models"
)

func (s *Store) GetCredential(id string) *models.Credential {
	s.mu.RLock()
	log.Println("Rlocked")
	defer func() {
		s.mu.RUnlock()
		log.Println("Runlocked")
	}()
	return s.credentials[id]
}

func (s *Store) GetCredentialsList() []*models.Credential {
	list := make([]*models.Credential, 0, len(s.credentials))
	s.mu.RLock()
	log.Println("Rlocked")
	defer func() {
		s.mu.RUnlock()
		log.Println("Runlocked")
	}()
	for _, v := range s.credentials {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
	return list
}
