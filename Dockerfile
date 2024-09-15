FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .
# run builds
RUN go build -o store ./client/main.go
RUN go build -o server ./server/main.go


#setup runtime env
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
# copy over builds
COPY --from=builder /app/store /usr/local/bin/store
COPY --from=builder /app/server /usr/local/bin/server

# dir for file storage
RUN mkdir -p /data/filestore
# FOR NOW design sets up folder as where to store files
ENV FILESTORE_DIR=/data/filestore
EXPOSE 8080

# run the server 
ENTRYPOINT ["/usr/local/bin/server/main"]