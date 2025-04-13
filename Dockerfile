FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN swag init --parseDependency --parseInternal -o docs -g ./cmd/init/main.go

RUN go build -o bin/server ./cmd/init/main.go

CMD ["/app/bin/server"]