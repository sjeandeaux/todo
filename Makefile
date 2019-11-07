BUILD_VERSION=0.0.0
OWNER=sjeandeaux
REPO=todo
SRC_DIR=github.com/$(OWNER)/$(REPO)

BUILD_TIME=$(shell date +%Y-%m-%dT%H:%M:%S%z)
GIT_COMMIT?=$(shell git rev-parse --short HEAD 2> /dev/null || echo "UNKNOWN")
GIT_DIRTY?=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_DESCRIBE?=$(shell git describe --tags --always 2> /dev/null || echo "UNKNOWN")
BUILD_TIME?=$(shell date +"%Y-%m-%dT%H:%M:%S")

LDFLAGS=-ldflags "\
          -X $(SRC_DIR)/pkg/information.Version=$(BUILD_VERSION) \
          -X $(SRC_DIR)/pkg/information.BuildTime=$(BUILD_TIME) \
          -X $(SRC_DIR)/pkg/information.GitCommit=$(GIT_COMMIT) \
          -X $(SRC_DIR)/pkg/information.GitDirty=$(GIT_DIRTY) \
          -X $(SRC_DIR)/pkg/information.GitDescribe=$(GIT_DESCRIBE)"

PKGGOFILES=$(shell go list ./... | grep -v todo-grpc)

# https://gist.github.com/sjeandeaux/e804578f9fd68d7ba2a5d695bf14f0bc
help: ## prints help.
	@grep -hE '^[a-zA-Z_-]+.*?:.*?## .*$$' ${MAKEFILE_LIST} | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: tools
tools: ## download tools
	go get -u github.com/grpc-ecosystem/grpc-health-probe
	go get -u github.com/client9/misspell/cmd/misspell
	go get -u golang.org/x/lint/golint
	go get -u github.com/fzipp/gocyclo
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u gotest.tools/gotestsum


.PHONY: dependencies
dependencies: ##Download the dependencies
	go mod download

.PHONY: build
build: 	##Build the binary
	mkdir -p ./target
	CGO_ENABLED=0 go build $(LDFLAGS) -installsuffix 'static' -o ./target/todod ./cmd/todod/main.go

.PHONY: gocyclo
gocyclo: ## check cyclomatic
	@gocyclo .

.PHONY: fmt
fmt: ## go fmt
	@go fmt $(PKGGOFILES)

.PHONY: misspell
misspell: ## go fmt on packages
	@misspell $(PKGGOFILES)

.PHONY: vet
vet: ## go vet on packages
	@go vet $(PKGGOFILES)

.PHONY: lint
lint: ## go lint on packages
	@golint -set_exit_status=true ./...

.PHONY: test
test: clean fmt vet ## test
	gotestsum --junitfile target/test-results/unit-tests.xml -- --short -cpu=2 -p=2 -coverprofile=target/coverage.txt -covermode=atomic -v $(LDFLAGS) $(PKGGOFILES)

.PHONY: it-test
it-test: clean fmt vet ## test
	gotestsum --junitfile target/test-results/it-tests.xml  -- -cpu=2 -p=2 -coverprofile=target/coverage.txt -covermode=atomic -v $(LDFLAGS) $(PKGGOFILES)

cover-html: it-test ## show the coverage in HTML page
	go tool cover -html=target/coverage.txt

clean: ## clean the target folder
	@rm -fr target
	@mkdir -p target/test-results

generate: ## generate the go from protobuf
	protoc --go_out=plugins=grpc:. todo-grpc/*.proto

docker-compose-build: ## builds the application image with docker-compose.
	VERSION=$(BUILD_VERSION) BUILD_DATE=$(BUILD_TIME) docker-compose build

docker-compose-up: ## spawns the containers.
	VERSION=$(BUILD_VERSION) BUILD_DATE=$(BUILD_TIME) docker-compose up -d

docker-build:
	docker build --tag $(OWNER)/$(REPO):$(BUILD_VERSION) .

docker-push:
	docker push $(OWNER)/$(REPO):$(BUILD_VERSION)