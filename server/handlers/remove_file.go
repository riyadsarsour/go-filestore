package handlers

import (
	"fmt"
	"go-filestore/server/filestore"
	"log"
	"net/http"
)

func RemoveFile(writer http.ResponseWriter, req *http.Request, store *filestore.FileStore, fileName string) {
	log.Println("handling remove request")
	err := store.RemoveFile(fileName)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to remove file: %v", err), http.StatusNotFound)
		return
	}

	writer.Write([]byte("File successfully removed"))
}
