FROM golang:1.23.5

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o main ./cmd/api/main.go
EXPOSE 3025

CMD ["./main"]