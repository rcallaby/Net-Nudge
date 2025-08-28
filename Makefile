.PHONY: build test lint run once


build:
go build -o bin/scansched ./cmd/scansched


test:
go test ./... -count=1 -race -timeout=5m


run:
go run ./cmd/scansched --config configs/sample.yaml


once:
go run ./cmd/scansched --config configs/sample.yaml --once