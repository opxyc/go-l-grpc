# Learning gRPC in Go - The Basics

Reference: https://medium.com/pantomath/how-we-use-grpc-to-build-a-client-server-system-in-go-dd20045fa1c2

## Sections

- [#1](../../tree/basic-client-server) Create proto file, generate go code, write a server and client
- [#2](../../tree/authenticating-server) Secure the communication - Authenticating the server
- [#3](../../tree/identifying-client) Secure the communication - Identifying the Clients
- [#4](../../tree/add-rest) Add REST Endpoints

---

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

## Securing the Communication

### Authenticating the Server (via certificate)

#### Create [self-signed SSL certificate](https://en.wikipedia.org/wiki/Self-signed_certificate)

```sh
# generate CA certificate
openssl genrsa -out cert/CA/CA.key -des3 2048
openssl req -x509 -sha256 -new -nodes -days 365 -key cert/CA/CA.key -out cert/CA/CA.crt
```

Add the below details to `cert/server/server.ext`:

```
authorityKeyIdentifier = keyid,issuer
basicConstraints = CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
IP.1 = 127.0.0.1
```

```sh
# generate certificate for server
openssl genrsa -out cert/server/server.key -des3 2048
openssl req -new -key cert/server/server.key -out cert/server/server.csr
openssl x509 -req -in cert/server/server.csr -CA cert/CA/CA.crt -CAkey cert/CA/CA.key -CAcreateserial -days 365 -sha256 -extfile cert/server/server.ext -out cert/server/server.crt
# decrypt the server's key
openssl rsa -in cert/server/server.key -out cert/server/server.decrypted.key
```

Above commands will generates server key and certificate, which we can use to authenticate the server. (we can also sign the certificate ourselves without the CA)

Reference: https://www.section.io/engineering-education/how-to-get-ssl-https-for-localhost/

### Identifying the Client

Another interesting feature of the gRPC server is the ability to intercept a request from the client. The client can inject information on the transport layer. We can use that feature to identify our client, because the SSL implementation authenticates the server (via the certificate), but not the client (all our clients are using the same certificate).
So we'll update the client side to inject metadata on every call (like a login and password), and the server side to check these credentials for every incoming call.

## Add REST Endpoints

As of now, have a pretty neat server, client and protocol; serialized, encrypted and authenticated. But there is a important limit: your client needs to be gRPC compliant, that is be in the list of [supported platforms](https://grpc.io/about/#osp). To avoid that limit, we can open the server to a REST gateway, allowing REST clients to perform requests it. Luckily, there is a gRPC [protoc plugin](https://github.com/grpc-ecosystem/grpc-gateway) to generate a reverse-proxy server which translates a RESTful JSON API into gRPC.

Refer https://github.com/grpc-ecosystem/grpc-gateway for details on working with API gaateway.

### Generating Go gateway code

- Install required binaries

  ```sh
  $ go install \
      github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest \
      github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest \
      google.golang.org/protobuf/cmd/protoc-gen-go@latest \
      google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
      github.com/contiamo/grpc-gateway/protoc-gen-swagger@latest
  ```

- Grab [googleapis repository](https://github.com/googleapis/googleapis) and place it to some location say, `$GOPATH/src/`
- Generate the code

  ```sh
  protoc -I $(GOPATH)/src/googleapis \
        --go_out=. \
        --go_opt=paths=source_relative \
        --go-grpc_out=.  \
        --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out . \
        --grpc-gateway_opt=paths=source_relative \
        --proto_path=. \
        --swagger_out=logtostderr=true:api \
        api/api.proto
  ```

  The above code will generate go code for proto, grpc and gateway, and swagger specification.
