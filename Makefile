APP=app

.PHONY: run build fmt test migrate-up migrate-down migrate-reset seed

run:
	go run ./cmd/$(APP) serve

build:
	go build -o bin/microseed ./cmd/$(APP)

fmt:
	gofmt -s -w .

test:
	go test ./...

migrate-up:
	go run ./cmd/$(APP) migrate up

migrate-down:
	go run ./cmd/$(APP) migrate down --step 1

migrate-reset:
	go run ./cmd/$(APP) migrate reset

seed:
	go run ./cmd/$(APP) seed
