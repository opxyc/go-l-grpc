SERVER_OUT := "server/server.bin"
CLIENT_OUT := "client/client.bin"
PKG := "github.com/opxyc/go-l-grpc"
SERVER_PKG_BUILD := "${PKG}/server"
CLIENT_PKG_BUILD := "${PKG}/client"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

.PHONY: api build build_server build_client

api: 
	@protoc  --go_out=. \
        --go_opt=paths=source_relative \
        --go-grpc_out=.  \
        --go-grpc_opt=paths=source_relative api/api.proto

dep: ## Get the dependencies
	@go get -v -d ./...

build: build_server build_client

build_server: dep api ## Build the binary file for server
	@go build -o $(SERVER_OUT) $(SERVER_PKG_BUILD)

build_client: dep api ## Build the binary file for client
	@go build -o $(CLIENT_OUT) $(CLIENT_PKG_BUILD)

clean: ## Remove previous builds
	@rm $(SERVER_OUT) $(CLIENT_OUT)

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
