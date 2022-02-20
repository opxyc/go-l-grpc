package main

import (
	"log"
	"net"

	"github.com/opxyc/go-l-grpc/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// main starts a gRPC server and waits for connection
func main() {
	lnsr, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a server instance
	s := &api.Server{}

	// create TLS credentials
	creds, err := credentials.NewServerTLSFromFile("cert/server/server.crt", "cert/server/server.decrypted.key")
	if err != nil {
		log.Fatalf("could not load certificate: %v", err)
	}

	// create an array of gRPC options with the credentials
	optns := []grpc.ServerOption{grpc.Creds(creds)}

	// create a gRPC server object
	grpcServer := grpc.NewServer(optns...)

	// attach the Ping service to the server
	api.RegisterPingServer(grpcServer, s)

	if err := grpcServer.Serve(lnsr); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
