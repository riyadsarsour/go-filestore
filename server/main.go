package main

import (
	"go-filestore/server/filestore"
	"go-filestore/server/handlers"
	"log"
	"net/http"
	"strings"
)

func main() {
	// need to revisit how i want to store files,
	// for now to running http server and verifying client usage this suffices
	// point for follow up in design talk
	fileStore := filestore.NewFileStore("./filestore")

	//  defining endpoints
	http.HandleFunc("/add", func(writer http.ResponseWriter, req *http.Request) {
		handlers.UploadFiles(writer, req, fileStore)
	})

	http.HandleFunc("/update", func(writer http.ResponseWriter, req *http.Request) {
		handlers.UpdateFile(writer, req, fileStore)
	})

	http.HandleFunc("/list", func(writer http.ResponseWriter, req *http.Request) {
		handlers.ListFiles(writer, req, fileStore)
	})

	http.HandleFunc("/remove/", func(writer http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodDelete {
			// Extract file name from URL path
			fileName := strings.TrimPrefix(req.URL.Path, "/remove/")
			if fileName == "" {
				http.Error(writer, "Missing file name", http.StatusBadRequest)
				return
			}
			handlers.RemoveFile(writer, req, fileStore, fileName)
		} else {
			http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/wordcount", func(writer http.ResponseWriter, req *http.Request) {
		handlers.WordCount(writer, req, fileStore)
	})

	http.HandleFunc("/freq-words", func(writer http.ResponseWriter, req *http.Request) {
		handlers.FrequentWords(writer, req, fileStore)
	})

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
