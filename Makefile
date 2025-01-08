include makefiles/security.mk

GOFMT_FILES?=$$(find . -not -path "./vendor/*" -type f -name '*.go')
APP_NAME?=terraform-provider-gocd
APP_DIR?=$$(git rev-parse --show-toplevel)
SRC_PACKAGES=$(shell go list -mod=vendor ./... | grep -v "vendor" | grep -v "mocks")
#VERSION?=0.0.1
VERSION := $(shell git tag --sort=-creatordate | head -1 | sed  -n 's/v\([0-9]*\).\([0-9]*\).\([0-9]*\)/\1.\2.\3/p')
NEXT_VERSION := $(shell echo ${VERSION}| awk -F '.' '{print $$1 "." $$2 "." $$3 +1 }' )

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

.PHONY: help
help: ## Prints help (only for targets with comments)
	@grep -E '^[a-zA-Z0-9._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; \
{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

local/fmt: ## Lints all the go code in the application.
	@gofmt -w $(GOFMT_FILES)
	$(GOBIN)/goimports -w $(GOFMT_FILES)
	$(GOBIN)/gofumpt -l -w $(GOFMT_FILES)
	$(GOBIN)/gci write $(GOFMT_FILES) --skip-generated

local/check: local/fmt ## Loads all the dependencies to vendor directory
	go mod tidy
	go mod vendor

build/local: local/check ## Generates the artifact with the help of 'go build'
	go build -o $(APP_NAME)_v$(VERSION) -ldflags="-s -w"

local/run: build.local ## Generates the artifact and start the service in the current directory
	./${APP_NAME}

lint: ## Lint's application for errors, it is a linters aggregator (https://github.com/golangci/golangci-lint).
	if [ -z "${DEV}" ]; then golangci-lint run --color always ; else docker run --rm -v $(APP_DIR):/app -w /app golangci/golangci-lint:v1.46.2-alpine golangci-lint run --color always ; fi

report: ## Publishes the go-report of the appliction (uses go-reportcard)
	docker run --rm -v $(APP_DIR):/app -w /app basnik/goreportcard-cli:latest goreportcard-cli -v

dev.prerequisite.up: ## Sets up the development environment with all necessary components.
	$(APP_DIR)/scripts/prerequisite.sh

generate.mock: ## generates mocks for the selected source packages.
	@go generate ${SRC_PACKAGES}

test: ## runs test cases
	@time go test $(TEST_FILES) -mod=vendor -coverprofile cover.out && go tool cover -html=cover.out -o cover.html && open cover.html

copy/terraformrc: ## copies terraformrc to user's home directory that helps in local development of the provider.
	cp terraformrc.sample ${HOME}/terraformrc

create/newversion/tfregistry: build/local ## Sets up the local terraform registry with the version specified.
	@mkdir -p ~/terraform-providers/registry.terraform.io/hashicorp/gocd/$(VERSION)/darwin_arm64/

upload/newversion/provider: create/newversion/tfregistry ## Uploads the updated provider to local terraform registry.
	@rm -rf  ~/terraform-providers/registry.terraform.io/hashicorp/gocd/$(VERSION)/darwin_arm64/terraform-provider-gocd_v$(VERSION)
	@cp terraform-provider-gocd_v$(VERSION) ~/terraform-providers/registry.terraform.io/hashicorp/gocd/$(VERSION)/darwin_arm64/

generate/document:
	tfplugindocs generate --website-source-dir templates/ --website-temp-dir templates-latest --examples-dir examples

tflint:
	@terraform fmt -write=false -check=true -diff=true examples/

version:
	@echo $(NEXT_VERSION)

build:
	GORELEASER_CURRENT_TAG=$(VERSION) goreleaser build --rm-dist --snapshot

local/build: local/check ## Generates the artifact with the help of 'go build'
	GORELEASER_CURRENT_TAG=$(VERSION) BUILD_ENVIRONMENT=${BUILD_ENVIRONMENT} goreleaser build --clean

publish: local/check ## Builds and publishes the app
	GOVERSION=${GOVERSION} BUILD_ENVIRONMENT=${BUILD_ENVIRONMENT} PLUGIN_PATH=${APP_DIR} goreleaser release --clean

mock/publish: local/check ## Builds and mocks app release
	GOVERSION=${GOVERSION} BUILD_ENVIRONMENT=${BUILD_ENVIRONMENT} PLUGIN_PATH=${APP_DIR} goreleaser release --skip=publish --clean --snapshot
