package storage

import (
	"io"
	"os"
	"path/filepath"
)

type DiskStore struct {
	baseDir string
}

func NewDiskStore(baseDir string) (*DiskStore, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	return &DiskStore{baseDir: baseDir}, nil
}

func (d *DiskStore) Save(id string, data io.Reader) (int64, error) {
	path := filepath.Join(d.baseDir, id)
	f, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return io.Copy(f, data)
}

func (d *DiskStore) Get(id string) (*os.File, error) {
	return os.Open(filepath.Join(d.baseDir, id))
}

func (d *DiskStore) Delete(id string) error {
	return os.Remove(filepath.Join(d.baseDir, id))
}
