FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o bin/server ./cmd/init/main.go

RUN make gen-docs

CMD ["/app/bin/server"]