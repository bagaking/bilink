package service

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bagaking/bilink/internal/index"
	"github.com/bagaking/bilink/internal/output"
	"github.com/bagaking/bilink/internal/rename"
	"github.com/bagaking/bilink/internal/resolve"
)

type RenameInput struct {
	Roots        []string
	OldPath      string
	NewPath      string
	Move         bool
	Extensions   []string
	ResolveRules resolve.Rules
	Interactive  bool
}

type pendingWrite struct {
	path     string
	tempPath string
	content  string
	mode     os.FileMode
}

func RunRename(input RenameInput) (output.RenamePayload, error) {
	extensions := input.Extensions
	if len(extensions) == 0 {
		extensions = []string{".md", ".markdown", ".mdx"}
	}
	idx, err := ScanAndIndex(input.Roots, extensions)
	if err != nil {
		return output.RenamePayload{}, err
	}
	oldBase := baseName(input.OldPath)
	newBase := baseName(input.NewPath)
	if err := ensureUnique(oldBase, idx, input.ResolveRules, input.Interactive); err != nil {
		return output.RenamePayload{}, err
	}
	var writes []pendingWrite
	paths := make([]string, 0, len(idx.Outbound))
	for path := range idx.Outbound {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return output.RenamePayload{}, err
		}
		content := string(data)
		rewritten := rename.RewriteWiki(content, oldBase, newBase)
		oldTargetSourcePath := path
		newTargetSourcePath := path
		if input.Move && samePath(path, input.OldPath) {
			newTargetSourcePath = input.NewPath
		}
		oldRel, newRel, err := markdownTargets(oldTargetSourcePath, newTargetSourcePath, input.OldPath, input.NewPath)
		if err != nil {
			return output.RenamePayload{}, err
		}
		rewritten = rename.RewriteMarkdown(rewritten, oldRel, newRel)
		if rewritten == content {
			continue
		}
		info, err := os.Stat(path)
		if err != nil {
			return output.RenamePayload{}, err
		}
		mode := info.Mode().Perm()
		writePath := path
		if input.Move && samePath(path, input.OldPath) {
			writePath = input.NewPath
		}
		writes = append(writes, pendingWrite{path: writePath, content: rewritten, mode: mode})
	}
	for i := range writes {
		tempPath, err := prepareTempWrite(writes[i].path, writes[i].content, writes[i].mode)
		if err != nil {
			cleanupTempWrites(writes)
			return output.RenamePayload{}, err
		}
		writes[i].tempPath = tempPath
	}
	if input.Move {
		if err := os.Rename(input.OldPath, input.NewPath); err != nil {
			cleanupTempWrites(writes)
			return output.RenamePayload{}, err
		}
	}
	updated := make([]string, 0, len(writes))
	for _, write := range writes {
		if err := os.Rename(write.tempPath, write.path); err != nil {
			cleanupTempWrites(writes)
			return output.RenamePayload{}, err
		}
		write.tempPath = ""
		if err := os.Chmod(write.path, write.mode); err != nil {
			cleanupTempWrites(writes)
			return output.RenamePayload{}, err
		}
		updated = append(updated, write.path)
	}
	return output.RenamePayload{OldPath: input.OldPath, NewPath: input.NewPath, Updated: updated}, nil
}

func prepareTempWrite(path string, content string, mode os.FileMode) (string, error) {
	dir := filepath.Dir(path)
	file, err := os.CreateTemp(dir, ".bilink-rename-*")
	if err != nil {
		return "", err
	}
	tempPath := file.Name()
	if _, err := file.WriteString(content); err != nil {
		_ = file.Close()
		_ = os.Remove(tempPath)
		return "", err
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(tempPath)
		return "", err
	}
	if err := os.Chmod(tempPath, mode); err != nil {
		_ = os.Remove(tempPath)
		return "", err
	}
	return tempPath, nil
}

func cleanupTempWrites(writes []pendingWrite) {
	for _, write := range writes {
		if write.tempPath != "" {
			_ = os.Remove(write.tempPath)
		}
	}
}

func markdownTargets(oldSourcePath, newSourcePath, oldPath, newPath string) (string, string, error) {
	oldRel, err := filepath.Rel(filepath.Dir(oldSourcePath), oldPath)
	if err != nil {
		return "", "", err
	}
	newRel, err := filepath.Rel(filepath.Dir(newSourcePath), newPath)
	if err != nil {
		return "", "", err
	}
	return filepath.ToSlash(oldRel), filepath.ToSlash(newRel), nil
}

func samePath(a, b string) bool {
	aAbs, aErr := filepath.Abs(a)
	bAbs, bErr := filepath.Abs(b)
	if aErr != nil || bErr != nil {
		return filepath.Clean(a) == filepath.Clean(b)
	}
	return filepath.Clean(aAbs) == filepath.Clean(bAbs)
}

func ensureUnique(name string, idx index.Index, rules resolve.Rules, interactive bool) error {
	if interactive {
		return nil
	}
	normalized := resolve.NormalizeName(name, rules)
	seen := map[string][]string{}
	for path := range idx.Outbound {
		base := baseName(path)
		key := resolve.NormalizeName(base, rules)
		seen[key] = append(seen[key], path)
	}
	if len(seen[normalized]) > 1 {
		return fmt.Errorf("ambiguous rename target: %s", name)
	}
	return nil
}

func baseName(path string) string {
	name := filepath.Base(path)
	if i := strings.LastIndex(name, "."); i > 0 {
		return name[:i]
	}
	return name
}
