package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// to add "add" arg
func uploadFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	body := &strings.Builder{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("files", filePath)
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	writer.Close()

	resp, err := http.Post("http://localhost:8080/add", writer.FormDataContentType(), strings.NewReader(body.String()))
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", resp.Status)
	}

	fmt.Println("File uploaded successfully")
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: store add <file1> [file2] ...")
		return
	}

	if os.Args[1] != "add" {
		fmt.Println("Unknown command:", os.Args[1])
		fmt.Println("Usage: store add <file1> [file2] ...")
		return
	}

	for _, file := range os.Args[2:] {
		if err := uploadFile(file); err != nil {
			fmt.Printf("Error uploading %s: %v\n", file, err)
		}
	}
}
