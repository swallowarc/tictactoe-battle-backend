# basic parameters
NAME     := tictactoe-battle-background
VERSION  := v0.0.0
REVISION := $(shell git rev-parse --short HEAD)

# Go parameters
BINARY_NAME=tictactoe-battle-background
SRCS    := $(shell find . -type f -name '*.go')
DIST_DIRS := find * -type d -exec
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""
GOOS = "linux"
GOARCH = "amd64"
GOCMD = go
GOBUILD = GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
GOVET = $(GOCMD) vet
GOGENERATE = $(GOCMD) generate
GOINSTALL = $(GOCMD) install

# build parameters
export GOPRIVATE=github.com/swallowarc/*
GRPC_PORT ?= 50051
DOCKER_CMD = docker
DOCKER_BUILD = $(DOCKER_CMD) build
DOCKER_PUSH =
DOCKER_REGISTRY = swallowarc/tictactoe-battle-backend
DOCKER_USER ?= fake_user
DOCKER_PASS ?= fake_pass

# test parameters
MOCK_DIR=internal/tests/mocks/
REDIS_HOST_PORT?=localhost:6379

.PHONY: build setup-tools upgrade-grpc mock-clean mock-gen vet test docker/build docker/push
build:
	$(GOBUILD) -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/ -v ./...
setup-tools:
	$(GOINSTALL) github.com/golang/mock/mockgen@v1.5.0
upgrade-grpc:
	$(GOGET) -u github.com/swallowarc/tictactoe_battle_proto
	$(GOMOD) tidy
mock-clean:
	rm -Rf ./$(MOCK_DIR)
mock-gen: mock-clean
	$(GOGENERATE) ./internal/domains/...
	$(GOGENERATE) ./internal/usecases/interactors/...
	$(GOGENERATE) ./internal/usecases/ports/...
	$(GOGENERATE) ./internal/interface_adapters/gateways/...
vet:
	$(GOVET) ./cmd/tictactoe_battle_backend/...
test:
	$(GOTEST) -v ./...
docker/build:
	$(DOCKER_BUILD) -t $(DOCKER_REGISTRY) .
docker/push:
