TEST ?= $(shell go list ./... | grep -v -e vendor -e keys -e tmp)
VERSION = $(shell cat version)
REVISION = $(shell git describe --always)
ifeq ("$(shell uname)","Darwin")
GO ?= GO111MODULE=on go
else
GO ?= GO111MODULE=on /usr/local/go/bin/go
endif
INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m
BUILD=tmp/bin

ci: depsdev test lint ## Run test and more...
test: ## Run test
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing$(RESET)"
	$(GO) test -v $(TEST) -timeout=30s -parallel=4
	$(GO) test -race $(TEST)

build: ## Build
	$(GO) build -o $(BUILD)/panalysis

run : ## Run
	$(GO) run main.go cli.go misc/sample.conf
lint: ## Exec golint
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Linting$(RESET)"
	golint -min_confidence 1.1 -set_exit_status $(TEST)
depsdev: ## Installing dependencies for development
	$(GO) get github.com/golang/lint/golint
	$(GO) get -u github.com/tcnksm/ghr
