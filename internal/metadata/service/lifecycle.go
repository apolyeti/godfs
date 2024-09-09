package metadata_service

import (
	"encoding/gob"
	dc "github.com/apolyeti/godfs/internal/data_node/client"
	metadata "github.com/apolyeti/godfs/internal/metadata/genproto"
	"log"
	"os"
	"sync"
	"time"
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
	shutdownChan chan struct{}
}

// NewMetadataService creates a new MetadataService

func NewMetadataService() *MetadataService {
	m := &MetadataService{
		inodes:       make(map[string]*Inode),
		numDataNodes: 3,
		dataNodes: []string{
			"data_node_1:50051",
			"data_node_2:50052",
			"data_node_3:50053",
		},
	}
	m.initializeRootDirectory()
	err := m.LoadFromDisk()
	if err != nil {
		log.Printf("No previous metadata found, starting with empty state")
	}

	go m.startHeartbeatLoop()
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

func (m *MetadataService) SendHeartbeat() {
	for _, dataNode := range m.dataNodes {
		client := dc.NewClient(dataNode)
		err := client.SendHeartbeat()

		if err != nil {
			log.Printf("Error sending heartbeat to %v: %v. Consider replacing this node", dataNode, err)
		}

		log.Println("Heartbeat successfully sent to", dataNode)
	}
}

func (m *MetadataService) startHeartbeatLoop() {
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			m.SendHeartbeat()
		case <-m.shutdownChan:
			log.Println("Shutting down heartbeat loop")
			ticker.Stop()
			return
		}
	}
}
