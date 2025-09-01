PROTO_SRC := protos
GOPATH_DIR := $(shell go env GOPATH)

ifeq ($(OS),Windows_NT)
    PROTOC_GEN_GO := $(GOPATH_DIR)\\bin\\protoc-gen-go.exe
    PROTOC_GEN_GRPC := $(GOPATH_DIR)\\bin\\protoc-gen-go-grpc.exe
	BIN_EXT := .exe
else
    PROTOC_GEN_GO := $(GOPATH_DIR)/bin/protoc-gen-go
    PROTOC_GEN_GRPC := $(GOPATH_DIR)/bin/protoc-gen-go-grpc
	BIN_EXT :=
endif

generate:
	@echo "Generating Go proto files..."
	@protoc \
		--go_out=. --go_opt=paths=import \
		--go-grpc_out=. --go-grpc_opt=paths=import \
		$(PROTO_SRC)/*.proto
	@echo "Done."

build:
	@echo "Building Go Program..."
	@go build -o ./build/FreeliPods$(BIN_EXT)
	@echo "Done."

install:
	@echo "Checking for Go Protobuf plugins..."
ifeq ($(OS),Windows_NT)
		@if not exist "$(PROTOC_GEN_GO)" ( \
			echo "protoc-gen-go not found. Installing now..." && \
			go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
		)
		@if not exist "$(PROTOC_GEN_GRPC)" ( \
			echo "protoc-gen-go-grpc not found. Installing now..." && \
			go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
		)
else
		@echo "Checking for Go Protobuf plugins..."
		@if [ ! -x "$(PROTOC_GEN_GO)" ]; then \
			echo "protoc-gen-go not found. Installing now..."; \
			go install google.golang.org/protobuf/cmd/protoc-gen-go@latest; \
		fi
		@if [ ! -x "$(PROTOC_GEN_GRPC)" ]; then \
			echo "protoc-gen-go-grpc not found. Installing now..."; \
			go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest; \
		fi
endif
	@echo "Go Protobuf plugins check complete."

clean:
	@echo "Cleaning generated Go proto files..."
	rm -rf $(OUT_DIR)/*.pb.go
	@echo "Done."

all: install generate build

.PHONY: all generate install build clean
