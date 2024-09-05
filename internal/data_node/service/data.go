// package data_node is responsible for storing and retrieving file data chunks.
// This is meant to be run on multiple nodes to distribute the storage of file data.

package data_service

import (
	"context"
	p "github.com/apolyeti/godfs/internal/data_node/genproto"
	"os"
)

const chunkSize = 4 * 1024 * 1024
const chunkDir = ".storage/chunks/"

// DataNode represents a node that stores file data chunks.
type DataNode struct {
	p.UnimplementedDataNodeServiceServer
	// ID of the data node, generated by the metadata service to
	// Keep track of which data nodes have which chunks
	// This ID will be used by the metadata service to determine which data nodes to send read requests to
	ID string
	// Map of chunk IDs to chunk data
	// For context, the chunk ID for files would be stored in the metadata service
	// When the client wants to read a file, it would get the chunk IDs from the metadata service
	// From there, it would request the chunk data from the data nodes
	Chunks map[string][]byte
}

// NewDataNode creates a new DataNode
func NewDataNode(id string) *DataNode {
	return &DataNode{
		ID:     id,
		Chunks: make(map[string][]byte),
	}
}
func (d *DataNode) WriteChunk(
	ctx context.Context,
	req *p.WriteChunkRequest,
) (
	*p.WriteChunkResponse, error,
) {
	err := os.MkdirAll(chunkDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(chunkDir + req.ChunkId)
	if err != nil {
		return nil, err
	}

	_, err = file.Write(req.Data)
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	d.Chunks[req.ChunkId] = req.Data
	return &p.WriteChunkResponse{}, nil
}

func (d *DataNode) ReadChunk(
	ctx context.Context,
	req *p.ReadChunkRequest,
) (
	*p.ReadChunkResponse, error,
) {
	data, ok := d.Chunks[req.ChunkId]
	if !ok {
		return nil, os.ErrNotExist
	}

	return &p.ReadChunkResponse{Data: data}, nil
}

func (d *DataNode) DeleteChunk(
	ctx context.Context,
	req *p.DeleteChunkRequest,
) (
	*p.DeleteChunkResponse, error,
) {
	err := os.Remove(chunkDir + req.ChunkId)
	if err != nil {
		return nil, err
	}

	delete(d.Chunks, req.ChunkId)
	return &p.DeleteChunkResponse{}, nil
}
