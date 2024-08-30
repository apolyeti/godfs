package metadata_client

import (
	"context"
	metaService "github.com/apolyeti/godfs/internal/metadata_service"
	metaGrpc "github.com/apolyeti/godfs/internal/metadata_service/genproto"
)

type Client struct {
	metadataClient metaGrpc.MetadataServiceClient
	currentDir     string
}

func NewClient(metadataClient metaGrpc.MetadataServiceClient) *Client {
	return &Client{
		metadataClient: metadataClient,
		currentDir:     metaService.RootID,
	}
}

func (c *Client) ChangeDir(dir string) error {
	req := &metaGrpc.ChangeDirRequest{
		CurrentDirectoryId: c.currentDir,
		TargetDirectoryId:  dir,
	}

	res, err := c.metadataClient.ChangeDir(context.Background(), req)
	if err != nil {
		return err
	}

	c.currentDir = res.DirectoryId
	return nil
}

func (c *Client) CurrentDir() string {
	return c.currentDir
}

func (c *Client) CreateFile(ctx context.Context, name string) (*metaGrpc.CreateFileResponse, error) {
	req := &metaGrpc.CreateFileRequest{
		Parent: c.currentDir,
		Name:   name,
	}

	return c.metadataClient.CreateFile(ctx, req)
}

func (c *Client) Mkdir(ctx context.Context, name string) (*metaGrpc.CreateFileResponse, error) {
	req := &metaGrpc.CreateFileRequest{
		Parent: c.currentDir,
		Name:   name,
		IsDir:  true,
	}

	return c.metadataClient.CreateFile(ctx, req)
}

func (c *Client) ListDir(ctx context.Context) (*metaGrpc.ListDirResponse, error) {
	req := &metaGrpc.ListDirRequest{
		DirectoryId: c.currentDir,
	}

	return c.metadataClient.ListDir(ctx, req)
}
