
build:
	go build -o bin/server ./cmd/server/main.go

run: 
	go run ./cmd/main.go

clean:
	rm -rf bin/server

deps:
	go mod tidy

lint:
	golangci-lint run

.DEFAULT_GOAL = run