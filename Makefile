.PHONY: build build-migrate run migrate dev clean test

build:
	go build -o bin/server ./cmd/server

build-migrate:
	go build -o bin/migrate ./cmd/migrate

run: build
	./bin/server

migrate: build-migrate
	./bin/migrate

dev: migrate run

clean:
	rm -rf bin/

test:
	go test ./...

install-deps:
	go mod download
	go mod tidy
