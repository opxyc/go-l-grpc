# Learning gRPC in Go - The Basics
Reference: https://medium.com/pantomath/how-we-use-grpc-to-build-a-client-server-system-in-go-dd20045fa1c2

## Sections
- [#1](tree/basic-client-server) Create proto file, generate go code, write a server and client

## Generating go code from protoc:
Install `protoc` and `protoc-gen-go`
```sh
sudo apt install -y protobuf-compiler
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```
Generate code using `protoc`
```sh
# from api/
protoc  --go_out=. \
        --go_opt=paths=source_relative \
        --go-grpc_out=.  \
        --go-grpc_opt=paths=source_relative api.proto
```

References:
- https://grpc.io/docs/languages/go/quickstart/