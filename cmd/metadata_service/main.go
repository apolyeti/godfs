// Entry point for the metadata_service service

package main

import (
	"github.com/apolyeti/godfs/internal/metadata_service"
	"google.golang.org/grpc"
)

func main() {
	metaDataService := metadata_service.NewMetadataService()

	// Setup gRPC server

	server := grpc.NewServer()

	metadata.RegisterMetadataServiceServer(server, metaDataService)
}
