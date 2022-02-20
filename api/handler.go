package api

import (
	"context"
	"log"
)

// Server representing the gRPC server
type Server struct{}

// SayHello generates response to a Ping request
func (s *Server) SayHello(ctx context.Context, in *PingMessage) (*PingMessage, error) {
	log.Printf("Received message %s", in.Greeting)
	return &PingMessage{Greeting: "bar"}, nil
}

func (s *Server) mustEmbedUnimplementedPingServer() {}
