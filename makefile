
# ==========================================
# Attor - database profiles and development
# ==========================================

DB ?= local
DB_TARGET_FILE := .db-target
MIGRATIONS_PATH := ./migrations

.PHONY: changedb db-status db-check run build \
        migrate-up migrate-down migrate-version migrate-force migrate-drop \
        fmt vet test build-all check tidy create-owner-test

# ------------------------------------------
# Database profile selection
# Supported profiles: local, attor-auth-test
# Production is deliberately excluded.
# ------------------------------------------

changedb:
	@if [ "$(DB)" != "local" ] && [ "$(DB)" != "attor-auth-test" ]; then \
		echo "ERROR: Unsupported database profile: $(DB)"; \
		echo "Allowed profiles: local, attor-auth-test"; \
		echo "Production switching is disabled."; \
		exit 1; \
	fi
	@if [ ! -f ".env.$(DB)" ]; then \
		echo "ERROR: .env.$(DB) not found"; \
		exit 1; \
	fi
	@if ! grep -qE '^[[:space:]]*DB_URL=' ".env.$(DB)"; then \
		echo "ERROR: DB_URL missing in .env.$(DB)"; \
		exit 1; \
	fi
	@printf '%s\n' "$(DB)" > "$(DB_TARGET_FILE)"
	@echo "Selected database profile: $(DB)"

# ------------------------------------------
# Show selected profile without credentials
# ------------------------------------------

db-status:
	@set -eu; \
	target=$$(cat "$(DB_TARGET_FILE)" 2>/dev/null) || { \
		echo "No database profile selected."; \
		echo "Run: make changedb DB=local"; \
		exit 1; \
	}; \
	case "$$target" in \
		local|attor-auth-test) ;; \
		*) echo "ERROR: Unsupported profile in $(DB_TARGET_FILE)"; exit 1 ;; \
	esac; \
	test -f ".env.$$target" || { echo "ERROR: .env.$$target missing"; exit 1; }; \
	echo "Selected database profile: $$target"

# ------------------------------------------
# Run API using the selected safe profile
# ------------------------------------------

run:
	@set -eu; \
	target=$$(cat "$(DB_TARGET_FILE)" 2>/dev/null) || { echo "Select a DB profile first"; exit 1; }; \
	case "$$target" in local|attor-auth-test) ;; *) echo "Unsupported profile: $$target"; exit 1 ;; esac; \
	test -f ".env.$$target" || { echo "Missing .env.$$target"; exit 1; }; \
	set -a; . "./.env.$$target"; set +a; \
	test -n "$${DB_URL:-}" || { echo "DB_URL missing"; exit 1; }; \
	go run ./cmd/api

# ------------------------------------------
# Verify database connection and identity
# ------------------------------------------

db-check:
	@set -eu; \
	target=$$(cat "$(DB_TARGET_FILE)" 2>/dev/null) || { echo "Select a DB profile first"; exit 1; }; \
	case "$$target" in local|attor-auth-test) ;; *) echo "Unsupported profile: $$target"; exit 1 ;; esac; \
	test -f ".env.$$target" || { echo "Missing .env.$$target"; exit 1; }; \
	echo "Checking profile: $$target"; \
	bash -c 'set -a; source "$$1"; set +a; exec go run ./cmd/db-check' _ ".env.$$target"

# ------------------------------------------
# Build
# ------------------------------------------

build:
	go build -o bin/api ./cmd/api

build-all:
	go build ./...

# ------------------------------------------
# Migrations
# Disabled until migration history is reconciled.
# Migrations 000004 and 000005 were applied manually
# to the Neon test branch.
# ------------------------------------------

migrate-up:
	@echo "BLOCKED: reconcile migration history before applying migrations."
	@echo "No database changes were made."
	@exit 1

migrate-down:
	@echo "BLOCKED: rollback is disabled until reviewed."
	@echo "No database changes were made."
	@exit 1

migrate-force:
	@echo "BLOCKED: do not force migration versions without reviewing database state."
	@exit 1

migrate-drop:
	@echo "BLOCKED: dropping a database is disabled."
	@exit 1

migrate-version:
	@set -eu; \
	target=$$(cat "$(DB_TARGET_FILE)" 2>/dev/null) || { echo "Select a DB profile first"; exit 1; }; \
	case "$$target" in local|attor-auth-test) ;; *) echo "Unsupported profile: $$target"; exit 1 ;; esac; \
	test -f ".env.$$target" || { echo "Missing .env.$$target"; exit 1; }; \
	set -a; . "./.env.$$target"; set +a; \
	test -n "$${DB_URL:-}" || { echo "DB_URL missing"; exit 1; }; \
	migrate -path "$(MIGRATIONS_PATH)" -database "$$DB_URL" version

# ------------------------------------------
# Go quality checks
# ------------------------------------------

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

check: fmt vet test

# ------------------------------------------
# Dependencies
# ------------------------------------------

tidy:
	go mod tidy

# ------------------------------------------
# Test-only owner creation
# ------------------------------------------

create-owner-test:
	@set -eu; \
	target=$$(cat "$(DB_TARGET_FILE)" 2>/dev/null) || { echo "Select a DB profile first"; exit 1; }; \
	test "$$target" = "attor-auth-test" || { \
		echo "ERROR: Select attor-auth-test first using make changedb DB=attor-auth-test"; \
		exit 1; \
	}; \
	test -f .env.attor-auth-test || { echo "Missing .env.attor-auth-test"; exit 1; }; \
	bash -c 'set -a; source ./.env.attor-auth-test; set +a; \
	test -n "$${DB_URL:-}" || { echo "DB_URL missing"; exit 1; }; \
	test -n "$${INITIAL_OWNER_EMAIL:-}" || { echo "INITIAL_OWNER_EMAIL missing"; exit 1; }; \
	test -n "$${INITIAL_OWNER_PASSWORD:-}" || { echo "INITIAL_OWNER_PASSWORD missing"; exit 1; }; \
	test -n "$${INITIAL_STORE_NAME:-}" || { echo "INITIAL_STORE_NAME missing"; exit 1; }; \
	exec go run ./cmd/create-owner'
