.PHONY: build clean

ifeq ($(strip $(BUILD_DIR)),)
BUILD_DIR=../../build
endif

clean:
	rm -f $(BUILD_DIR)/device

build:
	CGO_ENABLED=1 go build $(LDFLAGS) -gcflags "all=-trimpath=${GOPATH}" -o $(BUILD_DIR)/device main.go
