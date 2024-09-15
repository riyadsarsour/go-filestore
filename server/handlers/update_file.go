package handlers

import (
	"fmt"
	"go-filestore/server/filestore"
	"io"
	"log"
	"net/http"
)

func UpdateFile(writer http.ResponseWriter, req *http.Request, store *filestore.FileStore) {
	if req.Method != http.MethodPut {
		http.Error(writer, "Invalid request method, expected PUT request", http.StatusMethodNotAllowed)
		log.Printf("Invalid request method: %s", req.Method)
		return
	}

	log.Println("Starting file update process")

	reader, err := req.MultipartReader()
	if err != nil {
		http.Error(writer, "Unable to read multipart form", http.StatusBadRequest)
		log.Printf("Error reading multipart form: %v\n", err)
		return
	}

	//  while loop to iterate
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			http.Error(writer, "Error reading next part", http.StatusInternalServerError)
			log.Printf("Error reading next part: %v\n", err)
			return
		}

		log.Printf("Updating file: %s", part.FileName())

		err = store.UpdateFile(part)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Failed to update file: %v\n", err), http.StatusInternalServerError)
			log.Printf("Failed to update file: %s, Error: %v\n", part.FileName(), err)
			return
		}
		log.Printf("File successfully updated: %s\n", part.FileName())
	}
	writer.Write([]byte("File successfully updated or created"))
	log.Println("File update process completed")
}
