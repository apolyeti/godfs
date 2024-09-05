package main

import (
	"flag"
	dataGrpc "github.com/apolyeti/godfs/internal/data_node/genproto"
	service "github.com/apolyeti/godfs/internal/data_node/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	port := flag.String("port", "50051", "Port to listen on")
	flag.Parse()

	// Start the data node
	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := service.NewDataNode("localhost:" + *port)

	grpcServer := grpc.NewServer()

	dataGrpc.RegisterDataNodeServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	log.Printf("Serving on :%s", *port)
}
