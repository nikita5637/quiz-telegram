LOCAL_BIN:=$(CURDIR)/bin
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
GOTESTSUM_BIN=$(LOCAL_BIN)/gotestsum
PROTOC_GEN_GO_BIN:=$(LOCAL_BIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC_BIN:=$(LOCAL_BIN)/protoc-gen-go-grpc
GOIMPORTS_BIN:=$(LOCAL_BIN)/goimports

.PHONY: .bin-deps
.bin-deps:
ifeq ($(wildcard $(GOIMPORTS_BIN)),)
	$(info Installing binary dependency goimports)
	GOBIN=$(LOCAL_BIN) go install golang.org/x/tools/cmd/goimports 
endif
ifeq ($(wildcard $(GOTESTSUM_BIN)),)
	$(info Installing binary dependency gotestsum)
	GOBIN=$(LOCAL_BIN) go install gotest.tools/gotestsum
endif
ifeq ($(wildcard $(PROTOC_GEN_GO_BIN)),)
	$(info Installing binary dependency protoc-gen-go)
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go
endif
ifeq ($(wildcard $(PROTOC_GEN_GO_GRPC_BIN)),)
	$(info Installing binary dependency protoc-gen-go-grpc)
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc 
endif

.PHONY: .install-lint
.install-lint:
ifeq ($(wildcard $(GOLANGCI_BIN)),)
	$(info Downloading golangci-lint)
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint
	GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
endif

.PHONY: build
build:
	go build -o telegram ./cmd/telegram

.PHONY: clean
clean:
	rm telegram

.PHONY: generate
generate: .bin-deps
	PATH="$(LOCAL_BIN):$(PATH)" && protoc \
	--go_out=./pkg/pb/telegram --go_opt=paths=source_relative \
	--go-grpc_out=./pkg/pb/telegram --go-grpc_opt=paths=source_relative \
	--proto_path=./api/telegram ./api/telegram/telegram.proto

	$(GOIMPORTS_BIN) -w ./

.PHONY: go-generate
go-generate:
	go generate ./...

.PHONY: lint
lint: .install-lint
	$(info Running lint...)
	$(GOLANGCI_BIN) run --config=.golangci.pipeline.yaml ./...

.PHONY: migrations
migrations:
	goose -dir migrations/ mysql "$(USERNAME):$(DATABASE_PASSWORD)@tcp($(DBADDR))/$(DBNAME)" up

.PHONY: run
run:
	go run ./cmd/telegram --config ./config.yaml

.PHONY: test
test:
	$(GOTESTSUM_BIN) --format pkgname -- -coverprofile=cover.out ./internal/... ./utils/...
	go tool cover -html=cover.out
