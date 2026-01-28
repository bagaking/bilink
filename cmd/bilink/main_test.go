package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestRun_MissingCommand(t *testing.T) {
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stderr = w

	code := run([]string{})

	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	os.Stderr = oldStderr

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("close reader: %v", err)
	}

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(string(out), "missing command") {
		t.Fatalf("expected error output, got %q", string(out))
	}
}
