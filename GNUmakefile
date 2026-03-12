GOPATH ?= $(HOME)/go

default: fmt lint install generate

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.3 run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

# Run all acceptance tests (requires FLASHDUTY_API_KEY)
testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./internal/provider/...

# Run a specific acceptance test
# Usage: make testacc-run TEST=TestAccTeamResource
testacc-run:
	TF_ACC=1 go test -v -timeout 30m -run $(TEST) ./internal/provider/...

# Test data sources only
testacc-datasources:
	TF_ACC=1 go test -v -timeout 30m -run 'TestAcc.*DataSource' ./internal/provider/...

# Test resources only
testacc-resources:
	TF_ACC=1 go test -v -timeout 30m -run 'TestAcc.*Resource' ./internal/provider/...

# Quick smoke test - tests basic team and channel operations
testacc-smoke:
	TF_ACC=1 go test -v -timeout 10m -run 'TestAccTeam|TestAccChannel' ./internal/provider/...

# Build and install provider for local testing
local: build
	mkdir -p $(GOPATH)/bin
	cp terraform-provider-flashduty $(GOPATH)/bin/

clean:
	rm -f terraform-provider-flashduty
	go clean

.PHONY: fmt lint test testacc testacc-run testacc-datasources testacc-resources testacc-smoke build install generate local clean
