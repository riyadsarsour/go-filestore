package handlers

import (
	"fmt"
	"go-filestore/server/filestore"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func WordCount(writer http.ResponseWriter, req *http.Request, store *filestore.FileStore) {
	log.Println("Counting words ....")
	if req.Method != http.MethodGet {
		http.Error(writer, "Invalid request method, expected GET", http.StatusMethodNotAllowed)
		return
	}

	wordCount := 0

	err := filepath.WalkDir(store.Dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, _ := os.ReadFile(path)
		wordCount += len(strings.Fields(string(data)))
		return nil
	})

	if err != nil {
		http.Error(writer, "Failed to calculate word count", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(writer, "Total words: %d\n", wordCount)
	log.Println("Word Count done")
}
