GIT_TAG := $(shell echo $(shell git describe --tags || git branch --show-current) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')
BUILD_DATE	:= $(shell date '+%Y-%m-%d')

###############################################################################
###                                Build flags                              ###
###############################################################################

LD_FLAGS = -X github.com/EscanBE/escan-request-redirector/constants.VERSION=$(GIT_TAG) \
            -X github.com/EscanBE/escan-request-redirector/constants.COMMIT_HASH=$(COMMIT) \
            -X github.com/EscanBE/escan-request-redirector/constants.BUILD_DATE=$(BUILD_DATE)

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

###############################################################################
###                                  Test                                   ###
###############################################################################

test: go.sum
	@echo "testing"
	@go test -v ./... -race -coverprofile=coverage.txt -covermode=atomic
.PHONY: test

###############################################################################
###                                  Build                                  ###
###############################################################################

build: go.sum
	@echo "building Escan Request Redirect binary..."
	@echo "Flags $(BUILD_FLAGS)"
	@go build -mod=readonly $(BUILD_FLAGS) -o build/esrrd ./cmd/esrrd
	@echo "Builded successfully"
.PHONY: build

###############################################################################
###                                 Install                                 ###
###############################################################################

install: go.sum
	@echo "Installing Escan Request Redirect binary..."
	@echo "Flags $(BUILD_FLAGS)"
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/esrrd
	@echo "Installed successfully"
.PHONY: install