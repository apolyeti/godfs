package data_client

import (
	ctx "context"
	dataGrpc "github.com/apolyeti/godfs/internal/data_node/genproto"
)

type Client struct {
	DataNodeClient dataGrpc.DataNodeServiceClient
}

func NewClient(dataNodeClient dataGrpc.DataNodeServiceClient) *Client {
	return &Client{
		DataNodeClient: dataNodeClient,
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
