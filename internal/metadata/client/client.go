package metadata_client

import (
	"context"
	gRpc "github.com/apolyeti/godfs/internal/metadata/genproto"
	metaService "github.com/apolyeti/godfs/internal/metadata/service"
)

type Client struct {
	metadataClient gRpc.MetadataServiceClient
	currentDir     string
	currentDirName string
}

func NewClient(metadataClient gRpc.MetadataServiceClient) *Client {
	return &Client{
		metadataClient: metadataClient,
		currentDir:     metaService.RootID,
		currentDirName: "/",
	}
}

func (c *Client) ChangeDir(dir string) error {
	req := &gRpc.ChangeDirRequest{
		CurrentDirectoryId: c.currentDir,
		TargetDirectoryId:  dir,
	}

	res, err := c.metadataClient.ChangeDir(context.Background(), req)
	if err != nil {
		return err
	}

	c.currentDir = res.DirectoryId
	c.currentDirName = res.DirectoryName
	return nil
}

func (c *Client) CurrentDir() string {
	return c.currentDirName
}

func (c *Client) CurrentDirId() string { return c.currentDir }

func (c *Client) CreateFile(ctx context.Context,
	name string,
) (
	*gRpc.CreateFileResponse, error,
) {
	req := &gRpc.CreateFileRequest{
		Parent: c.currentDir,
		Name:   name,
	}

	return c.metadataClient.CreateFile(ctx, req)
}

func (c *Client) Mkdir(ctx context.Context,
	name string,
) (
	*gRpc.CreateFileResponse, error,
) {
	req := &gRpc.CreateFileRequest{
		Parent: c.currentDir,
		Name:   name,
		IsDir:  true,
	}

	return c.metadataClient.CreateFile(ctx, req)
}

func (c *Client) ListDir(ctx context.Context) (*gRpc.ListDirResponse, error) {
	req := &gRpc.ListDirRequest{
		DirectoryId: c.currentDir,
	}

	return c.metadataClient.ListDir(ctx, req)
}
