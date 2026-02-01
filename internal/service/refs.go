package service

import "github.com/bagaking/bilink/internal/output"

type RefsInput struct {
	Roots      []string
	Target     string
	Extensions []string
}

func RunRefs(input RefsInput) (output.RefsPayload, error) {
	idx, err := ScanAndIndex(input.Roots, input.Extensions)
	if err != nil {
		return output.RefsPayload{}, err
	}
	return output.RefsPayload{
		Target:   input.Target,
		Outbound: toAnySlice(idx.Outbound[input.Target]),
		Inbound:  toAnySlice(idx.Inbound[input.Target]),
	}, nil
}

func toAnySlice[T any](items []T) []any {
	if items == nil {
		return nil
	}
	out := make([]any, 0, len(items))
	for _, item := range items {
		out = append(out, item)
	}
	return out
}
