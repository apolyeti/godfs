package metadata_client

import (
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
