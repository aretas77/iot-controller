.PHONY: all clean purge build build-ui build-device

export PWD ?= $(shell pwd)
export BUILD_DIR ?= $(PWD)/build
export DATE := $(shell date -u +%y%m%d)

GIT_COMMIT ?= $(shell git rev-list -1 HEAD)
BUILD_TIME ?= $(shell date -u +%Y-%m-%d@%H:%M:%S)
HOSTNAME ?= $(shell hostname)
LDFLAGS ?= -ldflags "-installsuffix 'static' -w -s -X main.GitCommit=$(GIT_COMMIT) -X main.Date=$(BUILD_TIME) -X main.Host=$(HOSTNAME)"


all: clean build build-ui

purge:
	rm -rf $(BUILD_DIR)/*

clean:
	make -C $(PWD)/cmd/device clean

build:
	@go build $(LDFLAGS) -gcflags "all=-trimpath=${GOPATH}" -o ./build/main cmd/main.go

build-device:
	make -C $(PWD)/cmd/device build
