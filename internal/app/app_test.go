package app

import "testing"

func TestRunUnknownCommand(t *testing.T) {
	err := Run([]string{"nope"})
	if err == nil {
		t.Fatalf("expected error")
	}
}
