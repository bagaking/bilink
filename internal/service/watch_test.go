package service

import "testing"

func TestRunWatchRequiresIndex(t *testing.T) {
	_, err := RunWatch(WatchInput{IndexPath: ".bilink/index.json"})
	if err == nil {
		t.Fatalf("expected error")
	}
}
