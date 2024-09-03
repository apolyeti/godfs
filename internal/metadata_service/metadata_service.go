/*
Package metadata_service provides the metadata service for the metadata server.
*/

package metadata_service

import (
	"context"
	"encoding/gob"
	"errors"
	metadata "github.com/apolyeti/godfs/internal/metadata_service/genproto"
	"log"
	"os"
	"sync"
)

// MetadataService has 2 fields
// inodes: map of string to Inode
// mu: RWMutex for concurrent access to inodes
type MetadataService struct {
	metadata.UnimplementedMetadataServiceServer
	inodes map[string]*Inode
	mu     sync.RWMutex
}

// NewMetadataService creates a new MetadataService

func NewMetadataService() *MetadataService {
	m := &MetadataService{
		inodes: make(map[string]*Inode),
	}
	m.initializeRootDirectory()
	err := m.LoadFromDisk()
	if err != nil {
		log.Printf("No previous metadata found, starting with empty state")
	}
	return m
}

// RootID is the ID of the root directory
const RootID = "root"

func (m *MetadataService) initializeRootDirectory() {
	if _, ok := m.inodes[RootID]; !ok {
		m.inodes[RootID] = &Inode{
			ID:               RootID,
			Name:             "/",
			IsDir:            true,
			DirectoryEntries: map[string]string{},
		}
	}
}

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

func (m *MetadataService) SaveToDisk() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	file, err := os.Create(".storage/metadata.gob")
	if err != nil {
		return err
	}

	enc := gob.NewEncoder(file)
	err = enc.Encode(m.inodes)
	if err != nil {
		return err
	}

	return nil
}

func (m *MetadataService) LoadFromDisk() error {
	file, err := os.Open(".storage/metadata.gob")
	if err != nil {
		return err
	}

	dec := gob.NewDecoder(file)
	err = dec.Decode(&m.inodes)
	if err != nil {
		return err
	}

	return nil
}

func (m *MetadataService) Shutdown() {
	err := m.SaveToDisk()
	if err != nil {
		log.Printf("Error saving metadata to disk: %v", err)
	}
}
