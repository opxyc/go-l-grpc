package main

import (
	"context"
	"log"

	"github.com/opxyc/go-l-grpc/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Authentication struct {
	Login    string
	Password string
}

// GetRequestMetadata gets the current request metadata
func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"login":    a.Login,
		"password": a.Password,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security
func (a *Authentication) RequireTransportSecurity() bool {
	return true
}

func main() {
	creds, err := credentials.NewClientTLSFromFile("cert/server/server.crt", "")
	if err != nil {
		log.Fatalf("could not load certificate: %v", err)
	}

	auth := &Authentication{
		Login:    "john",
		Password: "doe",
	}

	conn, err := grpc.Dial(":7777", grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(auth))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := api.NewPingClient(conn)

	response, err := c.SayHello(context.Background(), &api.PingMessage{Greeting: "Hi"})
	if err != nil {
		log.Fatalf("could not call SayHello: %v", err)
	}

	log.Printf("response from server: %v", response.Greeting)
}
