package main

import (
	"io"
	"os"
	"testing"
)

func TestMain_PrintsBootstrap(t *testing.T) {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w

	main()

	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	os.Stdout = oldStdout

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("close reader: %v", err)
	}

	got := string(out)
	want := "bilink: MVP bootstrap\n"
	if got != want {
		t.Fatalf("unexpected output: got %q want %q", got, want)
	}
}
