PROJECT_NAME := "fantasy"
PKG := "github.com/tasylab/$(PROJECT_NAME)"
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _bin.go)
ALLENV := $(shellenv)


.PHONY: all dep build clean test coverage lint

all: build

lint: ## Lint the files
	@golangci-lint run  --deadline 10m

test: ## Run unittests
	@go test -short ./...

race: ## Run data race detector
	@go test -race -short ./...

coverage: ## Generate global code coverage report
	@mkdir -p /tmp/cover
	@go test -covermode=count -coverprofile "/tmp/cover/coverage.cov" ./...
	@go tool cover -func=/tmp/cover/coverage.cov

dep: ## Get the dependencies
	@GOPRIVATE=github.com/tasylab go mod download
	@go mod verify

server: dep ## Build the binary file
	@rm -rf bin
	@mkdir -p bin
	@CGO_ENABLED=0 go build -a -installsuffix cgo \
	-o bin/$(PROJECT_NAME) ./cmd/server


run: build ## Run the binary file
	@./bin/$(PROJECT_NAME)

clean: ## Remove previous build
	@rm -f bin/$(PROJECT_NAME)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2 }'