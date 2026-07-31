# =========================
# Variables
# =========================

DB_URL=postgres://postgres:postgres@localhost:5432/attor?sslmode=disable

MIGRATIONS_PATH=./migrations

# =========================
# Run
# =========================

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

# =========================
# Database
# =========================

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

migrate-force:
	@read -p "Version: " version; \
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $$version

migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version

migrate-drop:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" drop

# =========================
# Go
# =========================

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

build-all:
	go build ./...

# =========================
# Quality
# =========================

check: fmt vet test

# =========================
# Dependencies
# =========================

tidy:
	go mod tidy