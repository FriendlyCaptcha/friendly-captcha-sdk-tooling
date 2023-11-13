FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go mod download

# Use the go run command to start the server
CMD ["go", "run", "main.go", "serve", "--port", "1090"]
