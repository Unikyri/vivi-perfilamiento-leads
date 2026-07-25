.PHONY: build test vet run datos limpiar

build:
	go build -o bin/servidor ./cmd/servidor
	go build -o bin/pipeline ./cmd/pipeline

test:
	go test ./... -count=1

vet:
	go vet ./...

run:
	go run ./cmd/servidor

datos:
	go run ./cmd/pipeline

limpiar:
	rm -rf bin/
