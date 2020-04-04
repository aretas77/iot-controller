.PHONY: build clean

clean:
	rm -f $(BUILD_DIR)/device

build:
	CGO_ENABLED=1 go build $(LDFLAGS) -gcflags "all=-trimpath=${GOPATH}" -o $(BUILD_DIR)/device main.go
