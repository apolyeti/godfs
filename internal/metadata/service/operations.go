package metadata_service

import (
	"context"
	"fmt"
	dc "github.com/apolyeti/godfs/internal/data_node/client"
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

	conn, err := grpc.NewClient(
		dataNode,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("Failed to close connection: %v", err)
		}
	}()

	client := dc.NewClient(pb.NewDataNodeServiceClient(conn))

	err = client.WriteChunk(chunkId, chunkData)

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

	for i := 0; ; i++ {
		chunkId := fmt.Sprintf("%s-%d", inode.ID, i)
		dataNode := m.dataNodes[i%len(m.dataNodes)]

		chunkData, err := retrieveChunkFromDataNode(chunkId, dataNode)

		if err != nil {
			return nil, err
		}

		if len(chunkData) == 0 {
			break
		}

		data = append(data, chunkData...)
	}

	return &metadata.ReadFileResponse{
		Data: data,
	}, nil
}

func retrieveChunkFromDataNode(chunkId string, dataNode string) ([]byte, error) {
	var conn *grpc.ClientConn

	conn, err := grpc.NewClient(
		dataNode,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("Failed to close connection: %v", err)
		}
	}()

	client := dc.NewClient(pb.NewDataNodeServiceClient(conn))

	data, err := client.ReadChunk(chunkId)

	if err != nil {
		return nil, err
	}

	return data, nil
}
