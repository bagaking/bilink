package service

import "testing"

func TestCommandKind(t *testing.T) {
	if CommandRefs.String() == "" {
		t.Fatalf("expected command string")
	}
}
