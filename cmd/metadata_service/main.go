// Entry point for the metadata_service service

package main

import (
	metadata "github.com/apolyeti/godfs/internal/metadata_service"
	service "github.com/apolyeti/godfs/internal/metadata_service/genproto"
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

	service.RegisterMetadataServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	log.Printf("Serving on :8080")
}
