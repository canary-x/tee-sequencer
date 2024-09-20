GO=go
GO_TARGETS=./cmd/... ./internal/...
GIT_HASH=$(shell git rev-parse --short HEAD)
BUF_VERSION=1.40.1
PROTOC_GEN_GO_VERSION=1.34.2
PROTOC_GEN_CONNECT_GO_VERSION=1.16.2
GOLINT_VERSION=1.55.2

.PHONY: help
help:
	@grep -E '^[a-zA-Z_\-\/]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## build sequencer
	@${GO} build -o sequencer ./cmd/sequencer

.PHONY: docker/build
docker/build: ## build sequencer docker image
	docker build -t com.github.canary-x.tee-sequencer:latest .

.PHONY: build/enclave
build/enclave: docker/build ## build nitro enclave, only works on an EC2 instance with the nitro cli
	nitro-cli build-enclave --docker-uri com.github.canary-x.tee-sequencer:latest --output-file sequencer.eif

.PHONY: proto
proto: proto/lint proto/gen ## lint and generate proto files

.PHONY: proto/lint
proto/lint: ## lint proto files
	@buf lint ./proto
	@echo "Proto files lint successful"

.PHONY: proto/gen
proto/gen: ## generate proto sources
	@buf generate --template ./buf.gen.yaml .
	@echo "Proto files generated"

.PHONY: deps
deps: proto/setup ## set up all dependencies to run these make commands
	${GO} install github.com/golangci/golangci-lint/cmd/golangci-lint@v${GOLINT_VERSION}

.PHONY: proto/setup
proto/setup: ## install proto generation dependencies
	buf --version | grep ${BUF_VERSION} || ${GO} install github.com/bufbuild/buf/cmd/buf@v${BUF_VERSION}
	protoc-gen-go --version | grep ${PROTOC_GEN_GO_VERSION} || ${GO} install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GEN_GO_VERSION}

.PHONY: start
start: ## start sequencer as a Nitro instance
	@./scripts/start.sh
