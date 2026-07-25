.PHONY: build test vet run datos limpiar db-up db-down db-reset db-validate ci-local

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

DATABASE_URL ?= postgres://vivi:vivi@localhost:5432/vivi?sslmode=disable

db-up:
	docker compose up -d --wait

db-down:
	docker compose down -v

db-reset: db-down db-up

db-validate:
	@docker compose exec -T postgres psql -U vivi -d vivi -f /dev/stdin < migrations/001_esquema_inicial.sql > /dev/null
	@COUNT=$$(docker compose exec -T postgres psql -U vivi -d vivi -t -c "SELECT count(*) FROM information_schema.tables WHERE table_schema='public';" | tr -d ' \r\n'); \
	if [ "$$COUNT" = "7" ]; then echo "OK: 7 tables found"; else echo "FAIL: expected 7 tables, got $$COUNT"; exit 1; fi

ci-local: vet test build
	@echo "OK: listo para push"

# ─── Frontend ───────────────────────────────────────────────────────────────────

front-instalar:
	cd web && npm ci

front-dev:
	cd web && npm run dev

front-build:
	cd web && npm run build

front-verificar:
	cd web && npm run verificar
