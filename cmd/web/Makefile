.PHONY: build clean

ifeq ($(strip $(BUILD_DIR)),)
BUILD_DIR=../../build
endif

clean:
	rm -f $(BUILD_DIR)/web

build:
	go build $(LDFLAGS) -gcflags "all=-trimpath=${GOPATH}" -o $(BUILD_DIR)/web main.go
