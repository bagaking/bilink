package fs

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func ScanRoots(roots []string, extensions []string) ([]string, error) {
	extensionsSet := map[string]struct{}{}
	for _, ext := range extensions {
		extensionsSet[strings.ToLower(ext)] = struct{}{}
	}
	var files []string
	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			ext := strings.ToLower(filepath.Ext(path))
			if _, ok := extensionsSet[ext]; ok {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}
