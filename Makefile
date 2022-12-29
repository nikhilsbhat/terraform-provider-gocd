
GOFMT_FILES?=$$(find . -not -path "./vendor/*" -type f -name '*.go')
APP_NAME?=terraform-provider-gocd
APP_DIR?=$$(git rev-parse --show-toplevel)
SRC_PACKAGES=$(shell go list -mod=vendor ./... | grep -v "vendor" | grep -v "mocks")
VERSION?=0.0.1

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

local.fmt: ## Lints all the go code in the application.
	@gofmt -w $(GOFMT_FILES)
	$(GOBIN)/gofumpt -l -w $(GOFMT_FILES)
	$(GOBIN)/goimports -w $(GOFMT_FILES)

local.check: local.fmt ## Loads all the dependencies to vendor directory
	go mod vendor
	go mod tidy

local.build: local.check ## Generates the artifact with the help of 'go build'
	go build -o $(APP_NAME)_v$(VERSION) -ldflags="-s -w"

local.push: local.build ## Pushes built artifact to the specified location

local.run: local.build ## Generates the artifact and start the service in the current directory
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

create.newversion.tfregistry: local.build ## Sets up the local terraform registry with the version specified.
	@mkdir -p ~/terraform-providers/registry.terraform.io/hashicorp/gocd/$(VERSION)/darwin_arm64/

upload.newversion.provider: create.newversion.tfregistry ## Uploads the updated provider to local terraform registry.
	@rm -rf  ~/terraform-providers/registry.terraform.io/hashicorp/gocd/$(VERSION)/darwin_arm64/terraform-provider-gocd_v$(VERSION)
	@cp terraform-provider-gocd_v$(VERSION) ~/terraform-providers/registry.terraform.io/hashicorp/gocd/$(VERSION)/darwin_arm64/