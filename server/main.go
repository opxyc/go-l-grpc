package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/opxyc/go-l-grpc/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type contextKey int

const (
	clientIDKey contextKey = iota
)

func authenticateClient(ctx context.Context, s *api.Server) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		clientLogin := strings.Join(md["login"], "")
		clientPassword := strings.Join(md["password"], "")

		if clientLogin != "john" || clientPassword != "doe" {
			return "", fmt.Errorf("invalid username/password")
		}

		log.Printf("authenticated client: %s", clientLogin)
		return "42", nil
	}

	return "", fmt.Errorf("missing credentials")
}

// unaryInterceptor calls authenticationClient with current context
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	s, ok := info.Server.(*api.Server)
	if !ok {
		return nil, fmt.Errorf("unable to cast server")
	}

	clientID, err := authenticateClient(ctx, s)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, clientIDKey, clientID)
	return handler(ctx, req)
}

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
	// add a new grpc.ServerOption to the array 👇️
	optns := []grpc.ServerOption{grpc.Creds(creds), grpc.UnaryInterceptor(unaryInterceptor)}

	// create a gRPC server object
	grpcServer := grpc.NewServer(optns...)

	// attach the Ping service to the server
	api.RegisterPingServer(grpcServer, s)

	if err := grpcServer.Serve(lnsr); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
