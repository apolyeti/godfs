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

	c := metadata.NewMetadataServiceClient(conn)

	// create requests
	file1Req := &metadata.CreateFileRequest{
		Name:  "file1.txt",
		IsDir: false,
	}

	file2Req := &metadata.CreateFileRequest{
		Name:  "file2",
		IsDir: false,
	}

	// create files
	file1Res, err := c.CreateFile(context.Background(), file1Req)
	if err != nil {
		log.Fatalf("CreateFile failed: %v", err)
	}
	log.Printf("CreateFile Response: %v", file1Res)

	file2Res, err := c.CreateFile(context.Background(), file2Req)
	if err != nil {
		log.Fatalf("CreateFile failed: %v", err)
	}
	log.Printf("CreateFile Response: %v", file2Res)

	dirReq := &metadata.CreateFileRequest{
		Name:  "dir",
		IsDir: true,
	}

	dirRes, err := c.CreateFile(context.Background(), dirReq)
	if err != nil {
		log.Fatalf("CreateFile failed: %v", err)
	}
	log.Printf("CreateFile Response: %v", dirRes)

	// add file into directory
	file1Req.Parent = dirRes.Inode
	file1Res, err = c.CreateFile(context.Background(), file1Req)
	if err != nil {
		log.Fatalf("CreateFile failed: %v", err)
	}

	log.Printf("CreateFile Response: %v", file1Res)

	// list directory
	listDirReq := &metadata.ListDirRequest{
		DirectoryId: dirRes.Inode,
	}

	listDirRes, err := c.ListDir(context.Background(), listDirReq)
	if err != nil {
		log.Fatalf("ListDir failed: %v", err)
	}

	log.Printf("ListDir Response: %v", listDirRes)

}
