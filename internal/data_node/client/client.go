package data_client

import (
	ctx "context"
	dataGrpc "github.com/apolyeti/godfs/internal/data_node/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type Client struct {
	DataNodeClient dataGrpc.DataNodeServiceClient
}

func NewClient(dataNode string) *Client {
	var conn *grpc.ClientConn

	conn, err := grpc.NewClient(
		dataNode,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := dataGrpc.NewDataNodeServiceClient(conn)

	return &Client{
		DataNodeClient: client,
	}
}

func (c *Client) WriteChunk(chunkId string, data []byte) error {
	req := &dataGrpc.WriteChunkRequest{
		ChunkId: chunkId,
		Data:    data,
	}

	_, err := c.DataNodeClient.WriteChunk(ctx.Background(),
		req)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) ReadChunk(chunkId string) ([]byte, error) {
	req := &dataGrpc.ReadChunkRequest{
		ChunkId: chunkId,
	}

	res, err := c.DataNodeClient.ReadChunk(ctx.Background(),
		req)

	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (c *Client) DeleteChunk(chunkId string) error {
	req := &dataGrpc.DeleteChunkRequest{
		ChunkId: chunkId,
	}

	_, err := c.DataNodeClient.DeleteChunk(ctx.Background(),
		req)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SendHeartbeat() error {
	req := &dataGrpc.HeartbeatRequest{}

	_, err := c.DataNodeClient.Heartbeat(ctx.Background(),
		req)

	if err != nil {
		return err
	}

	return nil
}
