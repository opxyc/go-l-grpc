SERVER_OUT := "server/server.bin"
CLIENT_OUT := "client/client.bin"
PKG := "github.com/opxyc/go-l-grpc"
SERVER_PKG_BUILD := "${PKG}/server"
CLIENT_PKG_BUILD := "${PKG}/client"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
API_REST_OUT := "api/api.pb.gw.go"

.PHONY: api build build_server build_client

api: ## generate go pb, grpc and gateway code
	@protoc -I $(GOPATH)/src/googleapis \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=.  \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out . \
		--grpc-gateway_opt=paths=source_relative \
		--proto_path=. \
		--swagger_out=logtostderr=true:api \
		api/api.proto

dep: ## Get the dependencies
	@go get -v -d ./...

build: build_server build_client

build_server:
	@go build -o $(SERVER_OUT) $(SERVER_PKG_BUILD)

build_client:
	@go build -o $(CLIENT_OUT) $(CLIENT_PKG_BUILD)

clean: ## Remove previous builds
	@rm $(SERVER_OUT) $(CLIENT_OUT)

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
