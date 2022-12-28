package storage

import (
	"os"
	"path"
)

type Storage interface {
	Load(fileName string) ([]byte, error)
	Save(fileName string, content []byte) error
	RemoveFile(fileName string) error
	RemoveAll() error
}

type fileStorage struct {
	path string
}

// Load function reads file from disk and returns it content.
func (s *fileStorage) Load(fileName string) ([]byte, error) {
	data, err := os.ReadFile(path.Join(s.path, fileName))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Save function create or update file in the directory on the disk.
func (s *fileStorage) Save(fileName string, content []byte) error {
	if err := os.WriteFile(path.Join(s.path, fileName), content, 0o600); err != nil {
		return err
	}
	return nil
}

// RemoveAll function remove all files from directory.
func (s *fileStorage) RemoveAll() error {
	dir, err := os.ReadDir(s.path)
	if err != nil {
		return err
	}
	for _, d := range dir {
		err = os.RemoveAll(path.Join([]string{s.path, d.Name()}...))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *fileStorage) RemoveFile(fileName string) error {
	if err := os.Remove(path.Join(s.path, fileName)); err != nil {
		return err
	}
	return nil
}

// NewFileStorage creates new instance of file storage.
func NewFileStorage(path string) Storage {
	return &fileStorage{
		path: path,
	}
}
