package handlers

import (
	"fmt"
	"go-filestore/server/filestore"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// struct to hold freq map
type WordFrequency struct {
	text  string
	count int
}

// TODO Explore design to trigger a recount of frequent words upon file upload/update/delete
// for now this only runs when the endpoint is called
func FrequentWords(writer http.ResponseWriter, req *http.Request, store *filestore.FileStore) {

	query := req.URL.Query()
	order := query.Get("order")
	limitStr := query.Get("limit")

	// check if limit passed or use default limit
	limit := 10
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			// for now force limit to default if error occurs
			// i have client side issue errors on wrong command pass
			limit = 10
		}
	}

	wordFreq := make(map[string]int)

	err := filepath.WalkDir(store.Dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, _ := os.ReadFile(path)
		words := strings.Fields(string(data))

		for _, word := range words {
			wordFreq[word]++
		}
		return nil
	})

	if err != nil {
		http.Error(writer, "Failed to calculate frequent words", http.StatusInternalServerError)
		return
	}

	sortedWords := sortByFrequency(wordFreq, order, limit)

	for _, word := range sortedWords {
		fmt.Fprintf(writer, "%s: %d\n", word.text, word.count)
	}
}

// TODO Design wise need to handle tie-breaker rules (potentially just sort by alphabet)
func sortByFrequency(wordFreq map[string]int, order string, limit int) []WordFrequency {

	// if no order given will default to dsc for now
	if order == "" {
		order = "dsc"
	}

	var wordList []WordFrequency
	for word, count := range wordFreq {
		wordList = append(wordList, WordFrequency{word, count})
	}
	// println("word list")
	// fmt.Println(wordList)
	if order == "asc" {
		sort.Slice(wordList, func(i, j int) bool {
			return wordList[i].count < wordList[j].count
		})
	} else {
		sort.Slice(wordList, func(i, j int) bool {
			return wordList[i].count > wordList[j].count
		})
	}

	if limit > 0 && limit < len(wordList) {
		return wordList[:limit]
	} else if limit >= len(wordList) {
		return wordList
	}
	return wordList
}
