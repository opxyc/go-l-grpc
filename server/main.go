package main

import (
	"log"
	"net"

	"github.com/opxyc/go-l-grpc/api"
	"google.golang.org/grpc"
)

// main starts a gRPC server and waits for connection
func main() {
	lnsr, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a server instance
	s := &api.Server{}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach the Ping service to the server
	api.RegisterPingServer(grpcServer, s)

	if err := grpcServer.Serve(lnsr); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
