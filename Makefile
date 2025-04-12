
build:
	go build -o bin/server ./cmd/server/init/main.go

run: 
	go run ./cmd/init/main.go

clean:
	rm -rf bin/server

deps:
	go mod tidy

lint:
	golangci-lint run

migrate:
	go run ./cmd/migrate/main.go

.DEFAULT_GOAL = run