package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunRefs(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.md")
	b := filepath.Join(dir, "b.md")
	if err := os.WriteFile(a, []byte("See [[b]]"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := os.WriteFile(b, []byte("Hi"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	payload, err := RunRefs(RefsInput{Roots: []string{dir}, Target: b, Extensions: []string{".md"}})
	if err != nil {
		t.Fatalf("refs: %v", err)
	}
	inbound, ok := payload.Inbound.([]any)
	if !ok {
		t.Fatalf("expected inbound slice")
	}
	if len(inbound) == 0 {
		t.Fatalf("expected inbound")
	}
}
