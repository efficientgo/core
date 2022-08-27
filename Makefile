include .bingo/Variables.mk

GO111MODULE       ?= on
export GO111MODULE

GOBIN ?= $(firstword $(subst :, ,${GOPATH}))/bin

# Tools.
GIT ?= $(shell which git)

# Support gsed on OSX (installed via brew), falling back to sed. On Linux
# systems gsed won't be installed, so will use sed as expected.
SED ?= $(shell which gsed 2>/dev/null || which sed)

define require_clean_work_tree
	@git update-index -q --ignore-submodules --refresh

    @if ! git diff-files --quiet --ignore-submodules --; then \
        echo >&2 "$1: you have unstaged changes."; \
        git diff-files --name-status -r --ignore-submodules -- >&2; \
        echo >&2 "Please commit or stash them."; \
        exit 1; \
    fi

    @if ! git diff-index --cached --quiet HEAD --ignore-submodules --; then \
        echo >&2 "$1: your index contains uncommitted changes."; \
        git diff-index --cached --name-status -r --ignore-submodules HEAD -- >&2; \
        echo >&2 "Please commit or stash them."; \
        exit 1; \
    fi

endef

help: ## Displays help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: all
all: format build

.PHONY: build
build: ## Build modules
	go test -run=nope ./...

.PHONY: deps
deps: ## Cleans up deps for all modules
	@echo ">> running deps tidy"
	go mod tidy

.PHONY: docs
docs: $(MDOX) ## Generates config snippets and doc formatting.
	@echo ">> generating docs $(PATH)"
	@$(MDOX) fmt -l *.md

.PHONY: format
format: ## Formats Go code.
format: $(GOIMPORTS)
	@echo ">> formatting module"
	@$(GOIMPORTS) -w .

.PHONY: test
test: ## Runs all Go unit tests.
	@echo ">> running tests with -race"
	go test -v -race ./...

.PHONY: check-git
check-git:
ifneq ($(GIT),)
	@test -x $(GIT) || (echo >&2 "No git executable binary found at $(GIT)."; exit 1)
else
	@echo >&2 "No git binary found."; exit 1
endif

# PROTIP:
# Add
#      --cpu-profile-path string   Path to CPU profile output file
#      --mem-profile-path string   Path to memory profile output file
# to debug big allocations during linting.
lint: ## Runs various static analysis against our code.
lint: $(FAILLINT) $(GOLANGCI_LINT) $(MISSPELL) $(COPYRIGHT) build format docs check-git deps
	$(call require_clean_work_tree,"detected not clean master before running lint - run make lint and commit changes.")
	@echo ">> verifying imported packages"
	$(FAILLINT) -paths "fmt.{Print,PrintfPrintln,Sprint}" -ignore-tests ./...
	@echo ">> examining all of the Go files"
	go vet -stdmethods=false ./...
	@echo ">> linting all of the Go files GOGC=${GOGC}"
	$(GOLANGCI_LINT) run
	@echo ">> detecting misspells"
	@find . -type f | grep -v vendor/ | grep -vE '\./\..*' | xargs $(MISSPELL) -error
	@echo ">> ensuring Copyright headers"
	@ $(COPYRIGHT) $(shell go list -f "{{.Dir}}" ./... | xargs -i find "{}" -name "*.go")
	$(call require_clean_work_tree,"detected files without copyright - run make lint and commit changes.")
