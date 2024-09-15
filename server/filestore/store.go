package filestore

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileStore struct {
	Dir string
}

func NewFileStore(directory string) *FileStore {
	return &FileStore{Dir: directory}
}

func (s *FileStore) SaveFilePart(part *multipart.Part) error {
	filePath := filepath.Join(s.Dir, part.FileName())

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return fmt.Errorf("file already exists")
	}

	// else save
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	// copy over content
	if _, err := io.Copy(out, part); err != nil {
		return fmt.Errorf("failed to save file content: %v", err)
	}

	return nil
}
