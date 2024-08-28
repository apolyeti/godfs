package main

import (
	"context"
	"fmt"
	client "github.com/apolyeti/godfs/internal/metadata_client"
	metadata "github.com/apolyeti/godfs/internal/metadata_service/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	var conn *grpc.ClientConn

	conn, err := grpc.NewClient(
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("Failed to close connection: %v", err)
		}
	}()

	// Initialize the client
	c := client.NewClient(metadata.NewMetadataServiceClient(conn))

	res, err := c.CreateFile(context.Background(), "file1")

	if err != nil {
		log.Fatalf("Error creating file: %v\n", err)
	}

	fmt.Printf("File created: %v\n", res)

	res2, err := c.ListDir(context.Background())

	if err != nil {
		log.Fatalf("Error listing directory: %v\n", err)
	}

	fmt.Printf("Directory contents: %v\n", res2)

	// Make a directory now.

	res3, err := c.Mkdir(context.Background(), "dir1")

	if err != nil {
		log.Fatalf("Error creating directory: %v\n", err)
	}

	fmt.Printf("Directory created: %v\n", res3)

	res4, err := c.CreateFile(context.Background(), "file2")

	if err != nil {
		log.Fatalf("Error creating file: %v\n", err)
	}

	fmt.Printf("File created: %v\n", res4)

	res5, err := c.ListDir(context.Background())

	if err != nil {
		log.Fatalf("Error listing directory: %v\n", err)
	}

	fmt.Printf("Directory contents: %v\n", res5)

	// Change directory to dir1
	err = c.ChangeDir("dir1")

	if err != nil {
		log.Fatalf("Error changing directory: %v\n", err)
	}

	res6, err := c.CreateFile(context.Background(), "file3")

	if err != nil {
		log.Fatalf("Error creating file: %v\n", err)
	}

	fmt.Printf("File created: %v\n", res6)

	res7, err := c.ListDir(context.Background())

	if err != nil {
		log.Fatalf("Error listing directory: %v\n", err)
	}

	fmt.Printf("Directory contents: %v\n", res7)

}
