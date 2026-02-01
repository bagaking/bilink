package app

import "testing"

func TestRunUnknownCommand(t *testing.T) {
	err := Run([]string{"nope"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestRunRefsMissingArg(t *testing.T) {
	if err := Run([]string{"refs"}); err == nil {
		t.Fatalf("expected error")
	}
}
