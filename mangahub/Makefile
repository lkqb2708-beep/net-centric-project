.PHONY: all build run migrate seed test clean proto

# ─── Config ────────────────────────────────────────────────────────────────
BINARY     = mangahub
MAIN       = ./cmd/server
DB_URL    ?= $(shell grep DATABASE_URL .env | cut -d= -f2-)

# ─── Default ───────────────────────────────────────────────────────────────
all: build

# ─── Build ─────────────────────────────────────────────────────────────────
build:
	go build -o bin/$(BINARY) $(MAIN)

build-all:
	go build -o bin/server     ./cmd/server
	go build -o bin/tcpserver  ./cmd/tcpserver
	go build -o bin/udpserver  ./cmd/udpserver
	go build -o bin/grpcserver ./cmd/grpcserver

# ─── Run ───────────────────────────────────────────────────────────────────
run:
	go run $(MAIN)

run-tcp:
	go run ./cmd/tcpserver

run-udp:
	go run ./cmd/udpserver

run-grpc:
	go run ./cmd/grpcserver

# ─── Database ──────────────────────────────────────────────────────────────
migrate:
	go run ./cmd/server -migrate

seed:
	go run ./cmd/server -seed

db-setup:
	psql -U postgres -c "CREATE USER mangahub WITH PASSWORD 'mangahub_password';"
	psql -U postgres -c "CREATE DATABASE mangahub_db OWNER mangahub;"
	psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE mangahub_db TO mangahub;"

# ─── Proto ─────────────────────────────────────────────────────────────────
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       proto/manga.proto

# ─── Test ──────────────────────────────────────────────────────────────────
test:
	go test ./... -v -race -timeout 60s

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

# ─── Tidy ──────────────────────────────────────────────────────────────────
tidy:
	go mod tidy

# ─── Clean ─────────────────────────────────────────────────────────────────
clean:
	rm -rf bin/ coverage.out

# ─── Dev (all servers) ─────────────────────────────────────────────────────
dev:
	@echo "Starting all MangaHub servers..."
	@start cmd /k "go run ./cmd/tcpserver"
	@start cmd /k "go run ./cmd/udpserver"
	@start cmd /k "go run ./cmd/grpcserver"
	go run ./cmd/server
