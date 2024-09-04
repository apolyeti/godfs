/*
Package metadata_service provides the metadata service for the metadata server.
*/

package metadata_service

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/apolyeti/godfs/internal/data_node/genproto"
	metadata "github.com/apolyeti/godfs/internal/metadata/service/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

// GetInode returns the inode with the given ID
func (m *MetadataService) getInode(id string) (*Inode, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	inode, ok := m.inodes[id]
	if !ok {
		return nil, ErrFileNotFound
	}
	return inode, nil
}

func (m *MetadataService) CreateInode(
	ctx context.Context,
	req *metadata.CreateFileRequest,
) (
	*metadata.CreateFileResponse,
	error,
) {
	log.Printf("CREATEINODE\t%v", req)
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.inodes[req.Name]; ok {
		return nil, ErrExists
	}
	m.inodes[req.Name] = NewInode(req.Name, req.IsDir)
	return &metadata.CreateFileResponse{
		Name:  req.Name,
		Inode: m.inodes[req.Name].ID,
	}, nil
}

func (m *MetadataService) GetInode(
	ctx context.Context,
	req *metadata.GetInodeRequest,
) (
	*metadata.Inode,
	error,
) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	inode, ok := m.inodes[req.Name]
	if !ok {
		return nil, ErrFileNotFound
	}

	return &metadata.Inode{
		Name:  inode.Name,
		Id:    inode.ID,
		IsDir: inode.IsDir,
	}, nil
}

func (m *MetadataService) CreateFile(
	ctx context.Context,
	req *metadata.CreateFileRequest,
) (
	*metadata.CreateFileResponse,
	error,
) {
	log.Printf("CREATEFILE\t%v", req)

	m.mu.Lock()
	defer m.mu.Unlock()

	parentId := req.Parent

	if parentId == "" {
		parentId = RootID
	}

	parentInode, ok := m.inodes[parentId]
	if !ok {
		return nil, errors.New("parent directory not found")
	}

	if !parentInode.IsDir {
		return nil, ErrNotDir
	}

	if _, exists := parentInode.DirectoryEntries[req.Name]; exists {
		return nil, ErrExists
	}

	inode := NewInode(req.Name, req.IsDir)

	parentInode.DirectoryEntries[req.Name] = inode.ID

	m.inodes[inode.ID] = inode

	inode.ParentID = parentId

	return &metadata.CreateFileResponse{
		Name:  req.Name,
		Inode: inode.ID,
	}, nil
}
func (m *MetadataService) GetFile(
	ctx context.Context,
	req *metadata.CreateFileRequest,
) (
	*metadata.CreateFileResponse,
	error,
) {
	log.Printf("GETFILE\t%v", req)

	m.mu.Lock()
	defer m.mu.Unlock()

	parentInode, ok := m.inodes[req.Parent]
	if !ok {
		return nil, ErrFileNotFound
	}
	if !parentInode.IsDir {
		return nil, ErrNotDir
	}

	inodeId, exists := parentInode.DirectoryEntries[req.Name]
	if exists {
		return nil, ErrExists
	}

	inode, ok := m.inodes[inodeId]
	if !ok {
		return nil, ErrFileNotFound
	}
	if inode.IsDir {
		return nil, ErrIsDir
	}

	return &metadata.CreateFileResponse{
		Name:  req.Name,
		Inode: inode.ID,
	}, nil
}

func (m *MetadataService) listDir(inode *Inode) ([]*metadata.Inode, error) {

	if !inode.IsDir {
		return nil, ErrNotDir
	}

	var inodes []*metadata.Inode
	for _, id := range inode.DirectoryEntries {
		inode, ok := m.inodes[id]
		if !ok {
			return nil, ErrFileNotFound
		}
		inodes = append(inodes, &metadata.Inode{
			Name:  inode.Name,
			Id:    inode.ID,
			IsDir: inode.IsDir,
		})
	}
	return inodes, nil
}

func (m *MetadataService) ListDir(
	ctx context.Context,
	req *metadata.ListDirRequest,
) (
	*metadata.ListDirResponse,
	error,
) {
	log.Printf("LISTDIR\t%v", req)

	m.mu.Lock()
	defer m.mu.Unlock()

	var dirInode *Inode
	var ok bool

	if req.DirectoryId != "" {
		// Look up by DirectoryID
		dirInode, ok = m.inodes[req.DirectoryId]
		if !ok {
			return nil, ErrFileNotFound
		}
	} else if req.DirectoryName != "" && req.ParentId != "" {
		// Look up by DirectoryName and ParentID
		parentInode, ok := m.inodes[req.ParentId]
		if !ok {
			return nil, ErrFileNotFound
		}
		dirInodeID, exists := parentInode.DirectoryEntries[req.DirectoryName]
		if !exists {
			return nil, ErrFileNotFound
		}
		dirInode, ok = m.inodes[dirInodeID]
		if !ok {
			return nil, ErrFileNotFound
		}
	} else {
		return nil, errors.New("directory identifier not provided")
	}

	inodes, err := m.listDir(dirInode)
	if err != nil {
		return nil, err
	}

	return &metadata.ListDirResponse{
		Entries: inodes,
	}, nil
}

func (m *MetadataService) ChangeDir(
	ctx context.Context,
	req *metadata.ChangeDirRequest,
) (
	*metadata.ChangeDirResponse,
	error,
) {
	log.Printf("CHANGEDIR\t%v", req)

	m.mu.Lock()
	defer m.mu.Unlock()

	currentInode, ok := m.inodes[req.CurrentDirectoryId]
	if !ok {
		return nil, ErrFileNotFound
	}

	switch req.TargetDirectoryId {
	case "..":
		if currentInode.ID != RootID {
			parentInode, ok := m.inodes[currentInode.ParentID]
			if !ok {
				return nil, ErrFileNotFound
			}
			currentInode = parentInode
		}
	case ".":
		// Do nothing
	case "":
		// Go to root
		currentInode = m.inodes[RootID]
	default:
		// Because the client provides the NAME of the target directory,
		// We need to look up the ID of the target directory, which is done
		// Using the DirectoryEntries map, which maps the name to the ID of the directory inode.

		targetDirectoryId, exists := currentInode.DirectoryEntries[req.TargetDirectoryId]
		if !exists {
			return nil, ErrFileNotFound
		}

		targetInode, ok := m.inodes[targetDirectoryId]
		if !ok {
			return nil, ErrFileNotFound
		}

		currentInode = targetInode
	}

	return &metadata.ChangeDirResponse{
		DirectoryId:   currentInode.ID,
		DirectoryName: currentInode.Name,
	}, nil
}

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

	client := pb.NewDataNodeServiceClient(conn)

	_, err = client.WriteChunk(context.Background(), &pb.WriteChunkRequest{
		ChunkId: chunkId,
		Data:    chunkData,
	})

	if err != nil {
		return err
	}

	return nil
}
