package service

import (
	"fmt"
	"os"
	"path/filepath"
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
	var updated []string
	for path := range idx.Outbound {
		data, err := os.ReadFile(path)
		if err != nil {
			return output.RenamePayload{}, err
		}
		content := string(data)
		content = rename.RewriteWiki(content, oldBase, newBase)
		content = rename.RewriteMarkdown(content, filepath.Base(input.OldPath), filepath.Base(input.NewPath))
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return output.RenamePayload{}, err
		}
		updated = append(updated, path)
	}
	if input.Move {
		if err := os.Rename(input.OldPath, input.NewPath); err != nil {
			return output.RenamePayload{}, err
		}
	}
	return output.RenamePayload{OldPath: input.OldPath, NewPath: input.NewPath, Updated: updated}, nil
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
