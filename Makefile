.PHONY: run build tidy dev lint clean docker-dev-up docker-dev-down docker-prod

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

tidy:
	go mod tidy
	go mod verify

dev:
	air -c .air.toml

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ tmp/
	
docker-dev-up:
	docker compose -f docker-compose.dev.yml up --build

docker-dev-down:
	docker compose -f docker-compose.dev.yml down

docker-prod:
	docker compose up --build
