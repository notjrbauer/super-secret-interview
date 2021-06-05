SHELL=/bin/bash
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

.PHONY: gencerts 

