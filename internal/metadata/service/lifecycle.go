package metadata_service

import (
	"encoding/gob"
	metadata "github.com/apolyeti/godfs/internal/metadata/service/genproto"
	"log"
	"os"
	"sync"
)

// MetadataService has 2 fields
// inodes: map of string to Inode
// mu: RWMutex for concurrent access to inodes
type MetadataService struct {
	metadata.UnimplementedMetadataServiceServer
	inodes       map[string]*Inode
	mu           sync.RWMutex
	dataNodes    []string
	numDataNodes int
}

// NewMetadataService creates a new MetadataService

func NewMetadataService() *MetadataService {
	m := &MetadataService{
		inodes:       make(map[string]*Inode),
		numDataNodes: 3,
		dataNodes: []string{
			"localhost:50051",
			"localhost:50052",
			"localhost:50053",
		},
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

func (m *MetadataService) Shutdown() {
	err := m.SaveToDisk()
	if err != nil {
		log.Printf("Error saving metadata to disk: %v", err)
	}
}
