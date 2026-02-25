FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod tidy

COPY . .

RUN go build -o server ./cmd/server

EXPOSE 8080

CMD ["./server"]