package service

import "testing"

func TestRunCheck(t *testing.T) {
	payload, err := RunCheck(CheckInput{Roots: []string{"."}, Extensions: []string{".md"}})
	if err == nil {
		_ = payload
	}
}

func TestRunCheckGroups(t *testing.T) {
	payload, err := RunCheck(CheckInput{Roots: []string{"."}, Extensions: []string{".md"}})
	if err != nil {
		t.Fatalf("check: %v", err)
	}
	_ = payload.ErrorGroups
}
