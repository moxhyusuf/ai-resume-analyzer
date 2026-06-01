.PHONY: run build tidy test dev lint race clean

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

tidy:
	go mod tidy
	go mod verify

test:
	go test ./... -v

race:
	go test -race ./...

dev:
	swag init -g cmd/api/main.go --output docs && air -c .air.toml

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ tmp/

docker-build:
	docker compose up --build -d && docker compose logs -f