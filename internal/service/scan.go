package service

import (
	"os"

	"github.com/bagaking/bilink/internal/fs"
	"github.com/bagaking/bilink/internal/index"
)

func ScanAndIndex(roots []string, extensions []string) (index.Index, error) {
	files, err := fs.ScanRoots(roots, extensions)
	if err != nil {
		return index.Index{}, err
	}
	inputs := make([]index.FileInput, 0, len(files))
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			return index.Index{}, err
		}
		inputs = append(inputs, index.FileInput{Path: path, Content: string(data)})
	}
	return index.Build(inputs), nil
}
