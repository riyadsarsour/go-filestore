package main

import (
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

/* FURTURE TODOS FOR STRCUTURE OF CLIENT
	-- eventually introduce constant for endpoint to trigger when moving away from local host

END OF DESIGN?STRUCTURE NOTES */

// "add" logic
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
		errorMessage, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s", errorMessage)
	}

	fmt.Println("File uploaded successfully")
	return nil
}

func listFiles() error {
	resp, err := http.Get("http://localhost:8080/list")

	if err != nil {
		return fmt.Errorf("failed to list files: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// else if all defnesive checks pass we SHOULD have valid list of files
	fmt.Println("Files Stored:")
	fmt.Println(string(body))
	return nil
}

func removeFile(fileName string) error {
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/remove/"+fileName, nil)

	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to remove file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMessage, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s", errorMessage)
	}

	fmt.Println("File removed successfully")
	return nil
}

func updateFile(filePath string) error {
	file, err := os.Open(filePath)

	if err != nil {
		return fmt.Errorf("failed to open file: %v\n", err)
	}

	defer file.Close()

	body := &strings.Builder{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	writer.Close()

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/update", strings.NewReader(body.String()))

	if err != nil {
		return fmt.Errorf("failed to create update request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("failed to update file: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMessage, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s", errorMessage)
	}

	fmt.Println("File updated successfully")
	return nil
}

func wordCount() error {
	resp, err := http.Get("http://localhost:8080/wordcount")
	if err != nil {
		return fmt.Errorf("failed to get word count: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned:\n%v", resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	return nil
}

func getFreqWords(limit int, order string) error {

	req, err := http.NewRequest("GET", "http://localhost:8080/freq-words", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	query := req.URL.Query()
	if limit > 0 {
		query.Add("limit", strconv.Itoa(limit))
	}
	if order != "" {
		query.Add("order", order)
	}
	// eh right practice
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get freq-words: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	fmt.Println("Frequent Words:")
	fmt.Println(string(body))

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: store <command> [options]")
		fmt.Println("Commands:")
		fmt.Println(" add <file1> [file2] ...")
		fmt.Println(" update <file1>")
		fmt.Println(" wc")
		fmt.Printf(" freq-words [--limit|-n 10] [--order=dsc|")
		fmt.Println(" ls")
		fmt.Println(" remove <file>")
		fmt.Println(" -h | --help")
		return
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: store add <file1> [file2] ...")
			return
		}

		for _, file := range os.Args[2:] {
			if err := uploadFile(file); err != nil {
				fmt.Printf("Error uploading: %v\n", file)
				fmt.Println(err)
			}
		}
	case "update":
		if len(os.Args) < 3 {
			fmt.Println("Usage: store update <file1>")
			return
		}

		for _, file := range os.Args[2:] {
			if err := updateFile(file); err != nil {
				fmt.Printf("Error updating: %v\n", file)
				fmt.Println(err)
			}
		}
	case "ls":
		if len(os.Args) > 2 {
			fmt.Println("Too many args..\nUsage: store ls")
			return
		}

		if err := listFiles(); err != nil {
			fmt.Printf("Error listing files: %v\n", err)
		}
	case "remove":
		if len(os.Args) < 3 {
			fmt.Println("Too little Args\nUsage: store remove <file>")
			return
		}
		if err := removeFile(os.Args[2]); err != nil {
			fmt.Printf("Error removing file: %v\n", err)
		}
	case "wc":
		if len(os.Args) > 2 {
			fmt.Println("Too many args..\nUsage: store wc")
			return
		}
		if err := wordCount(); err != nil {
			fmt.Printf("Error getting word count: %v\n", err)
		}
	case "freq-words":
		limit := flag.Int("limit", 10, "Limit for the number of frequent words")
		flag.IntVar(limit, "n", 10, "Limit for the number of frequent words (short form)")

		order := flag.String("order", "dsc", "Order for frequent words, 'asc' or 'dsc' (default: dsc)")

		flag.CommandLine.Parse(os.Args[2:])

		if *limit <= 0 {
			fmt.Println("Error: Invalid value for --limit or -n. It must be a positive integer.")
			return
		}

		if *order != "asc" && *order != "dsc" {
			fmt.Println("Error: Invalid value for --order. Use 'asc' or 'dsc'.")
			return
		}

		if err := getFreqWords(*limit, *order); err != nil {
			fmt.Printf("Error fetching frequent words: %v\n", err)
		}

	case "-h", "--help":
		fmt.Println("Usage: store <command> [options]")
		fmt.Println("Commands:")
		fmt.Println("  add <file1> [file2] ...")
		fmt.Println("  update <file1>")
		fmt.Println("  wc")
		fmt.Println("  ls")
		fmt.Println("  freq-words [--limit|-n 10] [--order=dsc|")
		fmt.Println("  remove <file>")
		fmt.Println("  -h | --help")
	default:
		fmt.Println("Unknown command:", os.Args[1])
		fmt.Println("Usage: store <command> [options]")
		fmt.Println("Commands:")
		fmt.Println("  add <file1> [file2] ...")
		fmt.Println("  update <file1>")
		fmt.Println("  wc")
		fmt.Println("  freq-words [--limit|-n 10] [--order=dsc|asc]")
		fmt.Println("  ls")
		fmt.Println("  remove <file>")
		fmt.Println("  -h | --help")
	}
}
