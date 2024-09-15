package filestore

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func (s *FileStore) RemoveFile(fileName string) error {
	log.Printf("Removal triggered for file: %s", fileName)
	filePath := filepath.Join(s.Dir, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", fileName)
	}

	// else we have the file
	err := os.Remove(filePath)

	if err != nil {
		return fmt.Errorf("failed to remove file: %v", err)
	}

	log.Printf("File %s removed", fileName)
	return nil
}
