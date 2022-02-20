package main

import (
	"context"
	"log"

	"github.com/opxyc/go-l-grpc/api"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
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
