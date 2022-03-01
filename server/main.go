package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opxyc/go-l-grpc/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type contextKey int

const (
	clientIDKey contextKey = iota
)

func credMatcher(headerName string) (mdName string, ok bool) {
	if headerName == "Login" || headerName == "Password" {
		return headerName, true
	}
	return "", false
}

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

func startGRPCServer(address, certFile, keyFile string) error {
	lnsr, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// create a server instance
	s := &api.Server{}

	// create TLS credentials
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("could not load certificate: %v", err)
	}

	// create an array of gRPC options with the credentials
	// add a new grpc.ServerOption to the array 👇️
	optns := []grpc.ServerOption{grpc.Creds(creds), grpc.UnaryInterceptor(unaryInterceptor)}

	// create a gRPC server object
	grpcServer := grpc.NewServer(optns...)

	// attach the Ping service to the server
	api.RegisterPingServer(grpcServer, s)

	log.Printf("starting HTTP/2 gRPC server on %s", address)
	if err := grpcServer.Serve(lnsr); err != nil {
		return fmt.Errorf("failed to server: %v", err)
	}

	return nil
}

func startRESTServer(address, grpcAddress, certFile string) error {
	ctx := context.Background()
	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(credMatcher))

	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		return fmt.Errorf("could not load TLS certificate: %v", err)
	}

	// set up client gRPC options
	optns := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	// Register ping
	err = api.RegisterPingHandlerFromEndpoint(ctx, mux, grpcAddress, optns)
	if err != nil {
		return fmt.Errorf("could not register service 'Ping': %v", err)
	}

	log.Printf("starting HTTP/1.1 server on %s", address)
	http.ListenAndServe(address, mux)

	return nil
}

// main starts a gRPC server and waits for connection
func main() {
	grpcAddress := fmt.Sprintf("%s:%d", "localhost", 7777)
	restAddress := fmt.Sprintf("%s:%d", "localhost", 7778)
	certFile := "cert/server/server.crt"
	keyFile := "cert/server/server.decrypted.key"

	go func() {
		err := startGRPCServer(grpcAddress, certFile, keyFile)
		if err != nil {
			log.Fatalf("failed to start gRPC server: %v", err)
		}
	}()

	go func() {
		err := startRESTServer(restAddress, grpcAddress, certFile)
		if err != nil {
			log.Fatalf("failed to start gRPC server: %v", err)
		}
	}()

	select {}
}
