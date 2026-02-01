package output

import (
	"encoding/json"
	"fmt"
	"strings"
)

type RefsPayload struct {
	Target   string `json:"target"`
	Outbound any    `json:"outbound"`
	Inbound  any    `json:"inbound"`
}

type CheckPayload struct {
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

type RenamePayload struct {
	OldPath string   `json:"oldPath"`
	NewPath string   `json:"newPath"`
	Updated []string `json:"updated"`
}

type WatchPayload struct {
	Added   []string `json:"added"`
	Removed []string `json:"removed"`
}

func JSON(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func TextRefs(p RefsPayload) string {
	return fmt.Sprintf("target: %s\noutbound: %d\ninbound: %d", p.Target, count(p.Outbound), count(p.Inbound))
}

func TextCheck(p CheckPayload) string {
	return fmt.Sprintf("errors: %d\nwarnings: %d", len(p.Errors), len(p.Warnings))
}

func TextRename(p RenamePayload) string {
	return fmt.Sprintf("old: %s\nnew: %s\nupdated: %s", p.OldPath, p.NewPath, strings.Join(p.Updated, ", "))
}

func TextWatch(p WatchPayload) string {
	return fmt.Sprintf("added: %d\nremoved: %d", len(p.Added), len(p.Removed))
}

func count(v any) int {
	switch t := v.(type) {
	case []any:
		return len(t)
	case []string:
		return len(t)
	default:
		return 0
	}
}
