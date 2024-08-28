package metadata_client

import "github.com/apolyeti/godfs/internal/metadata_service/genproto/github.com/apolyeti/godfs/metadata_service"

type Client struct {
	metadataClient metadata_service.MetadataServiceClient
	currentDir     string
}
