// Entry point for the metadata_service service

package main

import (
	metadata "github.com/apolyeti/godfs/internal/metadata_service"
	service "github.com/apolyeti/godfs/internal/metadata_service/genproto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := metadata.NewMetadataService()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		log.Println("Shutting down server...")
		s.Shutdown()
		err := lis.Close()
		if err != nil {
			log.Fatalf("Error closing listener: %v", err)
		}
	}()

	grpcServer := grpc.NewServer()

	service.RegisterMetadataServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	log.Printf("Serving on :8080")
}
