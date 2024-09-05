package metadata_service

import (
	"context"
	"fmt"
	pb "github.com/apolyeti/godfs/internal/data_node/genproto"
	metadata "github.com/apolyeti/godfs/internal/metadata/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func (m *MetadataService) WriteFile(
	ctx context.Context,
	req *metadata.WriteFileRequest,
) (
	*metadata.WriteFileResponse,
	error,
) {
	log.Printf("WRITEFILE\t%v", req)

	m.mu.Lock()
	defer m.mu.Unlock()

	currDir := req.CurrentDirectoryId

	if currDir == "" {
		currDir = RootID
	}

	inodeId, exists := m.inodes[currDir].DirectoryEntries[req.FileName]

	if !exists {
		return nil, ErrFileNotFound
	}

	inode, ok := m.inodes[inodeId]

	if !ok {
		return nil, ErrFileNotFound
	}

	if inode.IsDir {
		return nil, ErrIsDir
	}

	chunks := chunkFile(req.Data, 1024)

	for i, chunk := range chunks {
		chunkId := fmt.Sprintf("%s-%d", inode.ID, i)
		dataNode := m.dataNodes[i%len(m.dataNodes)]

		err := storeChunkOnDataNode(chunkId, chunk, dataNode)
		inode.AddChunk(chunkId)

		if err != nil {
			return nil, err
		}
	}

	inode.UpdateSize(int64(len(req.Data)))

	return &metadata.WriteFileResponse{
		FileName: req.FileName,
	}, nil
}

func chunkFile(data []byte, chunkSize int) [][]byte {
	var chunks [][]byte

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}
	return chunks
}

func storeChunkOnDataNode(chunkId string, chunkData []byte, dataNode string) error {
	var conn *grpc.ClientConn

	// Correct method to dial gRPC server
	conn, err := grpc.Dial(
		dataNode,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // or WithInsecure()
	)

	if err != nil {
		return err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("Failed to close connection: %v", err)
		}
	}()

	client := pb.NewDataNodeServiceClient(conn)

	// Prepare the request for writing chunk data
	req := &pb.WriteChunkRequest{
		ChunkId: chunkId,
		Data:    chunkData,
	}

	_, err = client.WriteChunk(context.Background(), req)

	if err != nil {
		return err
	}

	return nil
}

func (m *MetadataService) ReadFile(
	ctx context.Context,
	req *metadata.ReadFileRequest,
) (
	*metadata.ReadFileResponse,
	error,
) {
	log.Printf("READFILE\t%v", req)

	m.mu.Lock()
	defer m.mu.Unlock()

	currDir := req.CurrentDirectoryId

	if currDir == "" {
		currDir = RootID
	}

	inodeId, exists := m.inodes[currDir].DirectoryEntries[req.FileName]

	if !exists {
		return nil, ErrFileNotFound
	}

	inode, ok := m.inodes[inodeId]

	if !ok {
		return nil, ErrFileNotFound
	}

	if inode.IsDir {
		return nil, ErrIsDir
	}

	var data []byte

	// Loop through stored chunks for the file
	for i, chunkId := range inode.ChunkIDs {
		dataNode := m.dataNodes[i%len(m.dataNodes)]

		chunkData, err := retrieveChunkFromDataNode(chunkId, dataNode)

		if err != nil {
			return nil, err
		}

		data = append(data, chunkData...)
	}

	return &metadata.ReadFileResponse{
		FileName: req.FileName,
		Data:     data,
	}, nil
}

func retrieveChunkFromDataNode(chunkId string, dataNode string) ([]byte, error) {
	var conn *grpc.ClientConn

	// Correct method to dial gRPC server
	conn, err := grpc.NewClient(
		dataNode,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // or WithInsecure()
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("Failed to close connection: %v", err)
		}
	}()

	client := pb.NewDataNodeServiceClient(conn)

	// Prepare the request for reading chunk data
	req := &pb.ReadChunkRequest{
		ChunkId: chunkId,
	}

	// Perform the RPC call to read the chunk data
	resp, err := client.ReadChunk(context.Background(), req)

	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}
