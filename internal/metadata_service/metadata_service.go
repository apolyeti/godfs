/*
Package metadata_service provides the metadata service for the metadata server.
*/

package metadata_service

import (
	"sync"
)

// MetadataService has 2 fields
// inodes: map of string to Inode
// mu: RWMutex for concurrent access to inodes
type MetadataService struct {
	inodes map[string]*Inode
	mu     sync.RWMutex
}

// NewMetadataService creates a new MetadataService

func NewMetadataService() *MetadataService {
	return &MetadataService{
		inodes: make(map[string]*Inode),
	}
}

// GetInode returns the inode with the given ID
func (m *MetadataService) GetInode(id string) (*Inode, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	inode, ok := m.inodes[id]
	if !ok {
		return nil, ErrFileNotFound
	}
	return inode, nil
}
