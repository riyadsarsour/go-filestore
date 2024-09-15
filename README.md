# go-filestore

## Table of Contents
- [go-filestore](#go-filestore)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Components](#components)
    - [Server](#server)
    - [Client](#client)
    - [CURRENT TEMP STATE NOTES](#current-temp-state-notes)
  - [Getting Statrted](#getting-statrted)
    - [OPTION 1: **Docker**](#option-1-docker)
      - [default way to run image](#default-way-to-run-image)
      - [OPtionally Specify Storage folder](#optionally-specify-storage-folder)
      - [Add Files](#add-files)
      - [List Files](#list-files)
      - [Remove a File](#remove-a-file)
      - [Update a File](#update-a-file)
      - [Word Count](#word-count)
      - [Frequent Words:](#frequent-words)
        - [Options:](#options)
    - [OPTION 2: Locally Cloning the Repository](#option-2-locally-cloning-the-repository)
    - [Server](#server-1)
    - [CLI Client](#cli-client)
  - [Commands and Usages](#commands-and-usages)
    - [Localhost CLI](#localhost-cli)
      - [Add Files](#add-files-1)
      - [List Files](#list-files-1)
      - [Remove a File](#remove-a-file-1)
      - [Update a File](#update-a-file-1)
      - [Word Count](#word-count-1)
      - [Frequent Words:](#frequent-words-1)
        - [Options:](#options-1)
  - [Tradeoffs \& Performance TODO Revisit](#tradeoffs--performance-todo-revisit)

## Overview

A simple service that to manage plain-text files through an HTTP server and a CLI client. 
The service supports file operations such as storing, updating, deleting files, and performing analysis on the stored files.

## Components

### Server

The server exposes HTTP endpoints to interact with the file store. It handles requests for adding, listing, removing, updating files, and performing file operations.

### Client

The CLI client interacts with the server. It supports commands to add, list, remove, update files, and perform operations such as word count and finding the most or least frequent words.

### CURRENT TEMP STATE NOTES 

Currently all files save to filestore directory in repo, however this can be changed in sever/main.go file
```go
fileStore := filestore.NewFileStore("YOUR/DESIRED/PATH")
```
## Getting Statrted

### OPTION 1: **Docker**

1. build the Docker image, run the build command from the root of the project directory
```
docker build -t go-file-store .
```
2. You can run the image 2 ways 
#### default way to run image
```
docker run -d -p 8080:8080 --name go-file-store go-file-store
```

#### OPtionally Specify Storage folder
```
docker run -d -p 8080:8080 -e FILESTORE_DIR=/your/storage/location --name go-file-store go-file-store
```
3. Interacting ***WITHIN** Docker Image via CLI
   #### Add Files

   Saves specified files in the current path to filestore directory. 
   Fails if the file already exists in the server.

   ```bash
   store add file1.txt file2.txt
   ```

   #### List Files

   Lists files stored.
   ```bash 
   store ls
   ```

   #### Remove a File

   Removes Specified files
   ```bash
   store rm file.txt
   ```

   #### Update a File

   Updates contents of a file on the server with the local file or creates a new file in server if absent.
   ```bash
   store update file.txt
   ```

   #### Word Count

   Returns the total number of words in all files stored in the server.

   ```bash
   store wc
   ```

   #### Frequent Words: 

   Returns the most or least frequent words in all files stored in the server.

   ```bash
   store freq-words [--limit|-n 10] [--order=dsc|asc]
   ```
   ##### Options:

   - `--limit | -n <number>`: Limits the number of results (default: 10).
   - `--order <asc|dsc>`: Sorts the result in ascending or descending order (default: `dsc`).

   
4. Interacting with Service via REST API
   ***ADD FILES**
   ```bash
   ccurl -X POST http://localhost:8080/add \
     -F "files=@file1.txt" \
     -F "files=@file2.txt" \
   ```
   Sample Reponse: 
   ```
   Files successfully uploaded
   ```
   ***List Files**
   ```bash
   curl -X GET http://localhost:8080/list
   ```
   Sample Reponse (assuming files exist): 
   ```
   file1.txt
   file2.txt
   ```
   *** Remove a File**
   ```bash
   curl -X DELETE http://localhost:8080/remove/file1.txt
   ```
   *** Update a File**
   ```bash 
   curl -X PUT http://localhost:8080/update -F "file=@file.txt"
   ```
   *** Word Count**
   ```bash
   curl -X GET http://localhost:8080/wordcount
   ```

   *** Frequent Words**
   ```bash
   curl -X GET "http://localhost:8080/freq-words?n=10&order=dsc"
   ```


### OPTION 2: Locally Cloning the Repository
### Server

CD to the server directory and build the server:

```bash
cd server
go build -o file-store-server
./file-store-server
```
The server will start on port 8080

### CLI Client

CD to the client directory and build the client

```bash
cd ../client
go build -o store
```
## Commands and Usages 

### Localhost CLI

#### Add Files

Saves specified files in the current path to filestore directory. 
Fails if the file already exists in the server.

```bash
./store add file1.txt file2.txt
```

#### List Files

Lists files stored.
```bash 
./store ls
```

#### Remove a File

Removes Specified files
```bash
./store rm file.txt
```

#### Update a File

Updates contents of a file on the server with the local file or creates a new file in server if absent.
```bash
./store update file.txt
```

#### Word Count

Returns the total number of words in all files stored in the server.

```bash
./store wc
```

#### Frequent Words: 

Returns the most or least frequent words in all files stored in the server.

```bash
./store freq-words [--limit|-n 10] [--order=dsc|asc]
```
##### Options:

- `--limit | -n <number>`: Limits the number of results (default: 10).
- `--order <asc|dsc>`: Sorts the result in ascending or descending order (default: `dsc`).


## Tradeoffs & Performance TODO Revisit

- **adding files**: Explore using hashing of content to see if new file's content already exists in another folder. Could help reduce redundant data transfers. Thouugh this will bring about complexity when handling file updates and this could be further explored with versioning.
- **Word Count**: The current implementation reads each file sequentially and counts words. This could be slow for large files, but can be optimized using parallel processing.
- **Frequent Words**: Sorting by frequency is handled on the server. This process could be memory-intensive, and future optimizing could include using a database or distributed file system.