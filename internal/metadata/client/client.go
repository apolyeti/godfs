package metadata_client

import (
	"context"
	"github.com/apolyeti/godfs/internal/metadata/genproto"
	metaService "github.com/apolyeti/godfs/internal/metadata/service"
)

type Client struct {
	metadataClient genproto.MetadataServiceClient
	currentDir     string
	currentDirName string
}

func NewClient(metadataClient genproto.MetadataServiceClient) *Client {
	return &Client{
		metadataClient: metadataClient,
		currentDir:     metaService.RootID,
		currentDirName: "/",
	}
}

func (c *Client) ChangeDir(dir string) error {
	req := &genproto.ChangeDirRequest{
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
	*genproto.CreateFileResponse, error,
) {
	req := &genproto.CreateFileRequest{
		Parent: c.currentDir,
		Name:   name,
	}

	return c.metadataClient.CreateFile(ctx, req)
}

func (c *Client) Mkdir(ctx context.Context,
	name string,
) (
	*genproto.CreateFileResponse, error,
) {
	req := &genproto.CreateFileRequest{
		Parent: c.currentDir,
		Name:   name,
		IsDir:  true,
	}

	return c.metadataClient.CreateFile(ctx, req)
}

func (c *Client) ListDir(ctx context.Context) (*genproto.ListDirResponse, error) {
	req := &genproto.ListDirRequest{
		DirectoryId: c.currentDir,
	}

	return c.metadataClient.ListDir(ctx, req)
}

func (c *Client) WriteFile(ctx context.Context, fileName string, data []byte) (*genproto.WriteFileResponse, error) {
	req := &genproto.WriteFileRequest{
		CurrentDirectoryId: c.currentDir,
		FileName:           fileName,
		Data:               data,
	}

	return c.metadataClient.WriteFile(ctx, req)
}

func (c *Client) ReadFile(ctx context.Context, fileName string) (*genproto.ReadFileResponse, error) {
	req := &genproto.ReadFileRequest{
		CurrentDirectoryId: c.currentDir,
		FileName:           fileName,
	}

	return c.metadataClient.ReadFile(ctx, req)
}
