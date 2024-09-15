package main

import (
	"go-filestore/server/filestore"
	"go-filestore/server/handlers"
	"log"
	"net/http"
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

	http.HandleFunc("/ls", func(writer http.ResponseWriter, req *http.Request) {
		handlers.ListFiles(writer, req, fileStore)
	})

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
