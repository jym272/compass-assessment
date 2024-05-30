# Makefile

all: format lint-fix
.PHONY: all

format:
	@gofmt -l -s -w .
.PHONY: format

lint:
	@golangci-lint run -c .golangci-gin.yml
.PHONY: lint

lint-fix:
	@golangci-lint run -c .golangci-gin.yml --fix
.PHONY: lint

test:
	@go test -v -cover ./...
.PHONY: test

build-app:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o score.app ./main.go
.PHONY: build-app
