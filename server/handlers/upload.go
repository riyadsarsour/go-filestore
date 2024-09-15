package handlers

import (
	"fmt"
	"go-filestore/server/filestore"
	"io"
	"log"
	"net/http"
)

func UploadFiles(writer http.ResponseWriter, req *http.Request, store *filestore.FileStore) {
	if req.Method != http.MethodPost {
		http.Error(writer, "Invalid request method, expected POST request", http.StatusMethodNotAllowed)
		log.Printf("Invalid request method: %s", req.Method)
		return
	}

	log.Println("Starting file upload process")

	reader, err := req.MultipartReader()
	if err != nil {
		http.Error(writer, "Unable to read multipart form", http.StatusBadRequest)
		log.Printf("Error reading multipart form: %v", err)
		return
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(writer, "Error reading next part", http.StatusInternalServerError)
			log.Printf("Error reading next part: %v", err)
			return
		}

		if part.FileName() == "" {
			continue
		}

		log.Printf("Processing file: %s", part.FileName())

		err = store.SaveFilePart(part)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Failed to save file: %v", err), http.StatusConflict)
			log.Printf("Failed to save file: %s, Error: %v", part.FileName(), err)
			return
		}
		log.Printf("File successfully saved: %s", part.FileName())
	}
	writer.Write([]byte("Files successfully uploaded"))
	log.Println("File upload process completed")

}
