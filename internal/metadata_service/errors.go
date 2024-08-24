// Package errors provides the custom errors used by the
// file system.

package metadata_service

import "errors"

var (
	ErrFileNotFound = errors.New("file not found")
	ErrDirNotFound  = errors.New("directory not found")
	ErrExists       = errors.New("file or directory already exists")
	ErrInvalidName  = errors.New("invalid file or directory name")
	ErrInvalidPath  = errors.New("invalid path")
	ErrIsDir        = errors.New("path is a directory")
	ErrNotEmpty     = errors.New("directory not empty")
	ErrNotDir       = errors.New("path is not a directory")
	ErrNotFile      = errors.New("path is not a file")
	ErrNotLink      = errors.New("path is not a link")
	ErrInvalidChunk = errors.New("invalid chunk")
	ErrInvalidSize  = errors.New("invalid size")
	ErrInvalidInode = errors.New("invalid inode")
)
