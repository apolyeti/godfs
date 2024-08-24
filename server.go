package main

import (
	"context"
	"github.com/apolyeti/godfs/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type ChatServer struct {
	chat.UnimplementedChatServer
}

func (s *ChatServer) Join(msg *chat.Message, stream chat.Chat_JoinServer) error {
	// Implement the logic for Join method
	return nil
}

// Send method implementation
func (s *ChatServer) Send(ctx context.Context, msg *chat.Message) (*chat.Message, error) {
	// Implement the logic for Send method
	return msg, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	chat.RegisterChatServer(s, &ChatServer{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}

	log.Println("gRPC server started on port :8080")
}
