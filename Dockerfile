FROM golang:1.20-alpine AS builder

WORKDIR /app
# Copy the Go modules manifest and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# build
RUN go build -o filestore ./server/main.go

FROM alpine:latest

# directory for the file store inside the container
# I could eventually do away with this but for now serves poc 
RUN mkdir -p /filestore

WORKDIR /app
# copy binary over
COPY --from=builder /app/filestore /app/filestore
EXPOSE 8080
# run the appl
CMD ["./filestore"]
