package main

import (
	"context"
	metadata "github.com/apolyeti/godfs/internal/metadata_service/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	var conn *grpc.ClientConn

	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := metadata.NewMetadataServiceClient(conn)

	res, err := c.CreateFile(context.Background(), &metadata.CreateFileRequest{
		Name:  "test",
		IsDir: false,
	})

	if err != nil {
		log.Fatalf("CreateFile failed: %v", err)
	}

	log.Printf("CreateFile Response: %v", res)
}
