
pkgs	:= $(shell go list ./...)

.PHONY: help lint test clean

ifdef GOPATH
GO111MODULE=on
endif

test: ## Run tests with coverage
	go test -short -cover -coverprofile coverage.out $(pkgs)
	go tool cover -html=coverage.out -o coverage.html

lint:
	golint $(pkgs)

clean: ## Clean all test files
	rm -rf coverage*

help: ## This help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) |sort |awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
