.PHONY: build clean

clean:
	rm -f $(BUILD_DIR)/web

build:
	go build $(LDFLAGS) -gcflags "all=-trimpath=${GOPATH}" -o $(BUILD_DIR)/web main.go
