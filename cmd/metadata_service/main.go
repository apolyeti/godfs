// Entry point for the metadata_service service

package main

import (
	metadata "github.com/apolyeti/godfs/internal/metadata_service"
	metadataservice "github.com/apolyeti/godfs/internal/metadata_service/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := metadata.NewMetadataService()

	grpcServer := grpc.NewServer()

	metadataservice.RegisterMetadataServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	log.Printf("Serving on :8080")
}
