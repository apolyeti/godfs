package main

import (
	"context"
	"fmt"
	client "github.com/apolyeti/godfs/internal/metadata/client"
	service "github.com/apolyeti/godfs/internal/metadata/genproto"
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
	c := client.NewClient(service.NewMetadataServiceClient(conn))

	res, err := c.CreateFile(context.Background(), "file1")

	if err != nil {
		log.Fatalf("Error creating file: %v\n", err)
	}

	fmt.Printf("File created: %v\n", res)

	res2, err := c.ListDir(context.Background())

	if err != nil {
		log.Fatalf("Error listing directory: %v\n", err)
	}

	fmt.Printf("Directory contents of %v:\n", c.CurrentDir())

	for _, entry := range res2.Entries {
		fmt.Printf("%v\n", entry.Name)
	}

	// Make a directory now.

	res3, err := c.Mkdir(context.Background(), "dir1")

	if err != nil {
		log.Fatalf("Error creating directory: %v\n", err)
	}

	fmt.Printf("Directory created: %v\n", res3.Name)

	res4, err := c.CreateFile(context.Background(), "file2")

	if err != nil {
		log.Fatalf("Error creating file: %v\n", err)
	}

	fmt.Printf("File created: %v\n", res4)

	res5, err := c.ListDir(context.Background())

	if err != nil {
		log.Fatalf("Error listing directory: %v\n", err)
	}

	fmt.Printf("Contents of %v:\n", c.CurrentDir())
	for _, entry := range res5.Entries {
		fmt.Printf("%v\n", entry.Name)
	}

	// Change directory to dir1
	err = c.ChangeDir("dir1")

	if err != nil {
		log.Fatalf("Error changing directory: %v\n", err)
	}

	res6, err := c.CreateFile(context.Background(), "file3")

	if err != nil {
		log.Fatalf("Error creating file: %v\n", err)
	}

	fmt.Printf("File created: %v\n", res6.Name)

	res7, err := c.ListDir(context.Background())

	if err != nil {
		log.Fatalf("Error listing directory: %v\n", err)
	}

	fmt.Printf("Directory contents of %v:\n", c.CurrentDir())

	for _, entry := range res7.Entries {
		fmt.Printf("%v\n", entry.Name)
	}

	fmt.Printf("Changing directory back to root\n")
	// Change directory back one
	err = c.ChangeDir("..")

	if err != nil {
		log.Fatalf("Error changing directory: %v\n", err)
	}

	res8, err := c.ListDir(context.Background())

	if err != nil {
		log.Fatalf("Error listing directory: %v\n", err)
	}

	fmt.Printf("Directory contents of %v:\n", c.CurrentDir())
	for _, entry := range res8.Entries {
		fmt.Printf("%v\n", entry.Name)
	}

	// Change directory to dir1
	err = c.ChangeDir("dir1")

	if err != nil {
		log.Fatalf("Error changing directory: %v\n", err)
	}

	res9, err := c.ListDir(context.Background())

	if err != nil {
		log.Fatalf("Error listing directory: %v\n", err)
	}

	fmt.Printf("Directory contents of %v:\n", c.CurrentDir())

	for _, entry := range res9.Entries {
		fmt.Printf("%v\n", entry.Name)
	}
}
