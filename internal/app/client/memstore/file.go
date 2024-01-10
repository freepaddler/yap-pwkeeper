package memstore

import (
	"errors"
	"os"
	"sort"

	"yap-pwkeeper/internal/pkg/models"
)

// GetFile returns File from store
func (s *Store) GetFileInfo(id string) *models.File {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.files[id]
}

// GetFilesList returns Files array from store sorted by Name
func (s *Store) GetFilesList() []*models.File {
	list := make([]*models.File, 0, len(s.files))
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.files {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
	return list
}

// AddFile stores File to server
func (s *Store) AddFile(d models.File, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return errors.New("failed to open file")
	}
	defer func() { _ = f.Close() }()
	st, err := os.Stat(filename)
	if err != nil {
		return errors.New("failed to stat file")
	}
	d.Size = st.Size()
	d.Filename = st.Name()
	return s.checkAuthErr(s.server.AddFile(d, f))
}

// UpdateFile Updates File on server
func (s *Store) UpdateFile(d models.File, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return errors.New("failed to open file")
	}
	defer func() { _ = f.Close() }()
	st, err := os.Stat(filename)
	if err != nil {
		return errors.New("failed to stat file")
	}
	d.Size = st.Size()
	d.Filename = st.Name()
	return s.checkAuthErr(s.server.UpdateFile(d, f))
}

// DeleteFile deletes File on server
func (s *Store) DeleteFile(d models.File) error {
	return s.checkAuthErr(s.server.DeleteFile(d))
}

// UpdateFileInfo updates FileInfo on server
func (s *Store) UpdateFileInfo(d models.File) error {
	return s.checkAuthErr(s.server.UpdateFileInfo(d))
}
