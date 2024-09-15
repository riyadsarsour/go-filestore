# go-filestore

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

1. **Docker**
   WORKING ON DOCKER IMAGE TO DOCKERHUB PLUS KUBERNETES DEPLOYMENT

2. **Clone the Repository**
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