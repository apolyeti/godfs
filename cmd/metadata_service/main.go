// Entry point for the metadata_service service

package main

import (
	metadata "github.com/apolyeti/godfs/internal/metadata_service"
	metadata_service "github.com/apolyeti/godfs/internal/metadata_service/service"
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

	metadata_service.RegisterMetadataServiceServer(grpcServer, s)

}
