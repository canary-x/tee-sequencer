GO=go
GO_TARGETS=./cmd/... ./internal/...
GIT_HASH=$(shell git rev-parse --short HEAD)

.PHONY: help
help:
	@grep -E '^[a-zA-Z_\-\/]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## build sequencer
	@${GO} build -o sequencer ./cmd/sequencer

.PHONY: docker/build
docker/build: ## build sequencer docker image
	docker build -t com.github.canary-x.tee-sequencer .
