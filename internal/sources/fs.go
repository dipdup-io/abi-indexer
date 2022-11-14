package sources

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileSystemConfig -
type FileSystemConfig struct {
	Dir string `yaml:"dir" validate:"required,dir"`
}

// FileSystem -
type FileSystem struct {
	root string
}

// NewFileSystem -
func NewFileSystem(rootDir string) *FileSystem {
	rootDir = filepath.Clean(rootDir)
	return &FileSystem{
		root: rootDir,
	}
}

// Get -
func (fs *FileSystem) Get(ctx context.Context, contract string) ([]byte, error) {
	filePath := filepath.Join(fs.root, fmt.Sprintf("%s.json", contract))
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return data, nil
}

// List -
func (fs *FileSystem) List(ctx context.Context) ([]string, error) {
	entries, err := os.ReadDir(fs.root)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for i := range entries {
		if entries[i].IsDir() {
			continue
		}

		info, err := entries[i].Info()
		if err != nil {
			return nil, err
		}
		parts := strings.Split(info.Name(), ".")
		result = append(result, parts[0])
	}

	return result, nil
}
