package memstore

import (
	"sort"

	"yap-pwkeeper/internal/pkg/models"
)

// GetCredential returns Credential from store
func (s *Store) GetCredential(id string) *models.Credential {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.credentials[id]
}

// GetCredentialsList returns Credentials array from store sorted by Name
func (s *Store) GetCredentialsList() []*models.Credential {
	list := make([]*models.Credential, 0, len(s.credentials))
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.credentials {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
	return list
}

// AddCredential saves new Credential to server
func (s *Store) AddCredential(d models.Credential) error {
	return s.checkAuthErr(s.server.AddCredential(d))
}

// UpdateCredential updates Credential on server
func (s *Store) UpdateCredential(d models.Credential) error {
	return s.server.UpdateCredential(d)
}

// DeleteCredential deletes Credential on server
func (s *Store) DeleteCredential(d models.Credential) error {
	return s.server.DeleteCredential(d)
}
