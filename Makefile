BIN := "./bin/image_previewer"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.48.0

lint: install-lint-deps
	golangci-lint run ./...

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/imagepreviewer/main.go

test:
	go test -race -count 100 ./...

integration-tests:
	docker-compose -f ./deployments/docker-compose.tests.yaml up --build --abort-on-container-exit --exit-code-from integration-tests && \
	docker-compose -f ./deployments/docker-compose.tests.yaml down

run:
	docker-compose -f ./deployments/docker-compose.yaml up -d

down:
	docker-compose -f ./deployments/docker-compose.yaml down

.PHONY: install-lint-deps lint build test integration-tests run down