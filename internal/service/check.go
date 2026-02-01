package service

import (
	"path/filepath"
	"strings"

	"github.com/bagaking/bilink/internal/check"
	"github.com/bagaking/bilink/internal/output"
	"github.com/bagaking/bilink/internal/resolve"
)

type CheckInput struct {
	Roots        []string
	Extensions   []string
	ResolveRules resolve.Rules
	LintRules    resolve.Rules
}

func RunCheck(input CheckInput) (output.CheckPayload, error) {
	idx, err := ScanAndIndex(input.Roots, input.Extensions)
	if err != nil {
		return output.CheckPayload{}, err
	}
	files := make([]check.File, 0, len(idx.Outbound))
	for path := range idx.Outbound {
		base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		files = append(files, check.File{
			Path:       path,
			ResolveKey: resolve.NormalizeName(base, input.ResolveRules),
			LintKey:    resolve.NormalizeName(base, input.LintRules),
		})
	}
	errs, warns := check.Detect(files)
	return output.CheckPayload{Errors: errs, Warnings: warns}, nil
}
