// Entry point for the metadata_service service

package main

import (
	p "github.com/apolyeti/godfs/internal/metadata/genproto"
	service "github.com/apolyeti/godfs/internal/metadata/service"
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

	s := service.NewMetadataService()

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

	p.RegisterMetadataServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	// occasionally log that the server is still running
	// and also send heartbeats tell metadata service to send heartbeats to data nodes

	log.Printf("Serving on :8080")
}
