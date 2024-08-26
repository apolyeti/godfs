/*
Package metadata_service provides the data structures used by the metadata service.

Inode represents a file or directory in the metadata service.
Ownership represents the ownership of a file or directory.
Timestamp represents the timestamps of a file or directory.
*/

package metadata_service

import (
	"github.com/google/uuid"
	"time"
)

// Ownership represents the ownership of a file or directory.
// UID: User ID of owner of file
// GID: Group ID of owner of file
type Ownership struct {
	UID int
	GID int
}

// Timestamp represents the timestamps of a file or directory.
// CreatedAt: Time at which file or directory was created
// UpdatedAt: Time at which file or directory was last updated
// AccessedAt: Time at which file or directory was last accessed
type Timestamp struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	AccessedAt time.Time
}

// Inode represents a file or directory in the metadata service.
// ID: Unique identifier of file or directory
// Name: Name of file or directory
// IsDir: True if inode is a directory, false if inode is a file
// Size: Size of file or directory in bytes
// Permissions: Permissions of file or directory
// Ownership: Ownership of file or directory
// Timestamp: Timestamps of file or directory
// ChunkIDs: IDs of chunks that store the data of the file
// ParentID: ID of parent directory
// Links: IDs of hard links to the file
type Inode struct {
	ID               string
	Name             string
	IsDir            bool
	Size             int64
	Permissions      string
	Ownership        Ownership
	Timestamp        Timestamp
	ChunkIDs         []string
	ParentID         string
	Links            []string
	DirectoryEntries []string
}

func NewInode(name string, isDir bool) *Inode {
	return &Inode{
		ID:          uuid.New().String(),
		Name:        name,
		IsDir:       isDir,
		Size:        0,
		Permissions: "rw-r--r--",
		Ownership:   Ownership{UID: 0, GID: 0},
		Timestamp: Timestamp{
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			AccessedAt: time.Now(),
		},
		ChunkIDs: []string{},
		ParentID: "",
		Links:    []string{},
	}
}

// AddLink adds a hard link to the inode.
func (i *Inode) AddLink(linkID string) {
	i.Links = append(i.Links, linkID)
}

// RemoveLink removes a hard link from the inode.
func (i *Inode) RemoveLink(linkID string) {
	for j, link := range i.Links {
		if link == linkID {
			i.Links = append(i.Links[:j], i.Links[j+1:]...)
			break
		}
	}
}

// AddChunk adds a chunk to the inode.
func (i *Inode) AddChunk(chunkID string) {
	i.ChunkIDs = append(i.ChunkIDs, chunkID)
}

// RemoveChunk removes a chunk from the inode.
func (i *Inode) RemoveChunk(chunkID string) {
	for j, chunk := range i.ChunkIDs {
		if chunk == chunkID {
			i.ChunkIDs = append(i.ChunkIDs[:j], i.ChunkIDs[j+1:]...)
			break
		}
	}
}

// UpdateSize updates the size of the inode.
func (i *Inode) UpdateSize(size int64) {
	i.Size = size
}

// UpdatePermissions updates the permissions of the inode.
func (i *Inode) UpdatePermissions(permissions string) {
	i.Permissions = permissions
}

// UpdateOwnership updates the ownership of the inode.
func (i *Inode) UpdateOwnership(ownership Ownership) {
	i.Ownership = ownership
}

// UpdateTimestamp updates the timestamps of the inode.
func (i *Inode) UpdateTimestamp(timestamp Timestamp) {
	i.Timestamp = timestamp
}

// UpdateParentID updates the parent ID of the inode.
func (i *Inode) UpdateParentID(parentID string) {
	i.ParentID = parentID
}

// UpdateName updates the name of the inode.
func (i *Inode) UpdateName(name string) {
	i.Name = name
}

// UpdateID updates the ID of the inode.
func (i *Inode) UpdateID(id string) {
	i.ID = id
}

// UpdateIsDir updates the isDir field of the inode.
func (i *Inode) UpdateIsDir(isDir bool) {
	i.IsDir = isDir
}

// UpdateLinks updates the links of the inode.
func (i *Inode) UpdateLinks(links []string) {
	i.Links = links
}

// UpdateChunkIDs updates the chunk IDs of the inode.
func (i *Inode) UpdateChunkIDs(chunkIDs []string) {
	i.ChunkIDs = chunkIDs
}

// GetLink returns the link at the specified index.
func (i *Inode) GetLink(index int) string {
	return i.Links[index]
}

// GetChunk returns the chunk at the specified index.
func (i *Inode) GetChunk(index int) string {
	return i.ChunkIDs[index]
}

// GetSize returns the size of the inode.
func (i *Inode) GetSize() int64 {
	return i.Size
}

// GetPermissions returns the permissions of the inode.
func (i *Inode) GetPermissions() string {
	return i.Permissions
}

// GetOwnership returns the ownership of the inode.
// Returns an Ownership struct
func (i *Inode) GetOwnership() Ownership {
	return i.Ownership
}

// GetTimestamp returns the timestamps of the inode.
// Returns a Timestamp struct
func (i *Inode) GetTimestamp() Timestamp {
	return i.Timestamp
}

// GetParentID returns the parent ID of the inode.
func (i *Inode) GetParentID() string {
	return i.ParentID
}

// GetName returns the name of the inode.
func (i *Inode) GetName() string {
	return i.Name
}

// GetID returns the ID of the inode.
func (i *Inode) GetID() string {
	return i.ID
}

// GetIsDir returns the isDir field of the inode.
func (i *Inode) GetIsDir() bool {
	return i.IsDir
}

// GetLinks returns the links of the inode.
func (i *Inode) GetLinks() []string {
	return i.Links
}

// GetChunkIDs returns the chunk IDs of the inode.
func (i *Inode) GetChunkIDs() []string {
	return i.ChunkIDs
}

// GetNumLinks returns the number of links to the inode.
func (i *Inode) GetNumLinks() int {
	return len(i.Links)
}

// GetNumChunks returns the number of chunks in the inode.
func (i *Inode) GetNumChunks() int {
	return len(i.ChunkIDs)
}
