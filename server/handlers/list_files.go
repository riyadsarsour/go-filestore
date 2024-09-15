package handlers

import (
	"fmt"
	"go-filestore/server/filestore"
	"net/http"
	"os"
)

func ListFiles(writer http.ResponseWriter, req *http.Request, store *filestore.FileStore) {
	if req.Method != http.MethodGet {
		http.Error(writer, "Invalid Request, 'ls' command expected GET request", http.StatusMethodNotAllowed)
		return
	}

	files, err := os.ReadDir(store.Dir)
	if err != nil {
		err_str := fmt.Sprintf("Unable to list files %e", err)
		http.Error(writer, err_str, http.StatusInternalServerError)
	}

	for _, file := range files {
		if file.IsDir() {
			// for now we are ignoring folder
			// need to revist logic to saving folders of work
			continue
		}

		writer.Write([]byte(file.Name() + "\n"))
	}
}
