SHELL=/bin/bash
BIN_DIR = $(PWD)/bin
GO_GEN_PATH="proto"
PACKAGES := $(shell go list ./... )
REPO=$(shell pwd)
NAME="worker"

proto:
	$(shell protoc --proto_path=$(REPO)/$(GO_GEN_PATH) --go_out=$(REPO)/$(GO_GEN_PATH) --go_opt=paths=source_relative \
		--go-grpc_out=$(GO_GEN_PATH) --go-grpc_opt=paths=source_relative \
		$(REPO)/proto/worker.proto; \
	)
	@echo generated $(shell find proto -name *.pb.go)

gencerts:
	bash scripts/gen-certs.sh

api: 
	go build -o $(BIN_DIR)/$@ cmd/$@/main.go 

cli: 
	go build -o $(BIN_DIR)/$@ cmd/$@/main.go 

test:
	go test $(REPO)/... -v -race -failfast -cover

build: api cli

clean: 
	rm -rf bin/*

all: clean build 

.PHONY: gencerts api cli test all

