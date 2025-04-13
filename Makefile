
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

gen-docs:
	swag init --parseDependency --parseInternal -o docs -g ./cmd/init/main.go

migrate:
	go run ./cmd/migrate/main.go

.DEFAULT_GOAL = run