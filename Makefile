BINARY_NAME=kasir-api
BUILD_DIR=bin
CMD_PATH=cmd/kasir-api/main.go

.PHONY: all build test coverage clean

all: build

build:
	mkdir -p $(BUILD_DIR)
	GCO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.txt ./...
	go tool cover -func=coverage.txt

clean:
	rm -rf $(BUILD_DIR) coverage.txt

lint:
	golangci-lint run

fmt:
	golangci-lint fmt
	
