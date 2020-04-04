export PWD ?= $(shell pwd)
export BUILD_DIR ?= $(PWD)/build
export DATE := $(shell date -u +%y%m%d)

# The python import module is located in the MQTT library
export PYTHONPATH = $(PWD)/../paho.mqtt.golang/

GIT_COMMIT ?= $(shell git rev-list -1 HEAD)
BUILD_TIME ?= $(shell date -u +%Y-%m-%d@%H:%M:%S)
HOSTNAME ?= $(shell hostname)
CMD ?=

export LDFLAGS ?= -ldflags "-installsuffix 'static' -w -s -X main.GitCommit=$(GIT_COMMIT) -X main.Date=$(BUILD_TIME) -X main.Host=$(HOSTNAME)"
export CGO_LDFLAGS = -g -O2 -L${GOPATH}/src/github.com/tensorflow/tensorflow/lite/tools/make/gen/linux_x86_64/lib/

# Default command is to start
ifeq ($(strip $(CMD)),)
CMD = start
endif

.PHONY: all clean purge build device docker

all: build

all-clean: clean purge

purge:
	rm -rf $(BUILD_DIR)/*

clean:
	@make -C $(PWD)/cmd/device clean
	@make -C $(PWD)/cmd/web clean

build:
	@make -C $(PWD)/cmd/web build
	@make -C $(PWD)/cmd/device build

docker:
	@docker-compose up --remove-orphans

device:
	PYTHONPATH=$(PYTHONPATH) ./build/device $(CMD)
