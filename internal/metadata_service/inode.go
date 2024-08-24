/*
Package metadata_service provides the data structures used by the metadata service.

Inode represents a file or directory in the metadata service.
Ownership represents the ownership of a file or directory.
Timestamp represents the timestamps of a file or directory.
*/

package metadata_service

import (
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
	ID          string
	Name        string
	IsDir       bool
	Size        int64
	Permissions string
	Ownership   Ownership
	Timestamp   Timestamp
	ChunkIDs    []string
	ParentID    string
	Links       []string
}

// TODO: Implement the following methods for Inode:
// - AddLink
// - RemoveLink
// - AddChunk
// - RemoveChunk
// - UpdateSize
// - UpdatePermissions
// - UpdateOwnership
// - UpdateTimestamp
// - UpdateParentID
// - UpdateName
// - UpdateID
// - UpdateIsDir
// - UpdateLinks
// - UpdateChunkIDs
// - GetLink
// - GetChunk
// - GetSize
// - GetPermissions
// - GetOwnership
// - GetTimestamp
// - GetParentID
// - GetName
// - GetID
// - GetIsDir
// - GetLinks
// - GetChunkIDs
// - GetNumLinks
// - GetNumChunks
