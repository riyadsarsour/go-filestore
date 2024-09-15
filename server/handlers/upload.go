package handlers

import (
	"fmt"
	"go-filestore/server/filestore"
	"io"
	"net/http"
)

func UploadFiles(writer http.ResponseWriter, req *http.Request, store *filestore.FileStore) {
	if req.Method != http.MethodPost {
		http.Error(writer, "Invalid request method, expected POST request", http.StatusMethodNotAllowed)
		return
	}

	reader, err := req.MultipartReader()
	if err != nil {
		http.Error(writer, "Unable to read multipart form", http.StatusBadRequest)
		return
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(writer, "Error reading next part", http.StatusInternalServerError)
			return
		}

		if part.FileName() == "" {
			continue
		}

		err = store.SaveFilePart(part)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Failed to save file: %v", err), http.StatusConflict)
			return
		}
	}

	writer.Write([]byte("Files successfully uploaded"))
}
