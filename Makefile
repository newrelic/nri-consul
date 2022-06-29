WORKDIR      := $(shell pwd)
NATIVEOS	 := $(shell go version | awk -F '[ /]' '{print $$4}')
NATIVEARCH	 := $(shell go version | awk -F '[ /]' '{print $$5}')
INTEGRATION  := consul
BINARY_NAME   = nri-$(INTEGRATION)
GO_FILES     := ./src/
GOFLAGS			= -mod=readonly


all: build

build: clean test compile integration-test

clean:
	@echo "=== $(INTEGRATION) === [ clean ]: Removing binaries and coverage file..."
	@rm -rfv bin coverage.xml

compile:
	@echo "=== $(INTEGRATION) === [ compile ]: Building $(BINARY_NAME)..."
	@go build -o bin/$(BINARY_NAME) $(GO_FILES)

test:
	@echo "=== $(INTEGRATION) === [ test ]: Running unit tests..."
	@go test -race ./... -count=1

integration-test:
	@echo "=== $(INTEGRATION) === [ test ]: running integration tests..."
	@go test -v -tags=integration ./test/integration/. -count=1

# Include thematic Makefiles
include $(CURDIR)/build/ci.mk
include $(CURDIR)/build/release.mk


.PHONY: all build clean compile test
