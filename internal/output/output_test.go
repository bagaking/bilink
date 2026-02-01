package output

import (
	"strings"
	"testing"
)

func TestJSONRefsOutput(t *testing.T) {
	payload := RefsPayload{Target: "a.md"}
	data, err := JSON(payload)
	if err != nil {
		t.Fatalf("json: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("expected json")
	}
}

func TestTextCheckIncludesGroups(t *testing.T) {
	payload := CheckPayload{
		Errors: []string{"foo"},
		ErrorGroups: []CheckGroup{{Key: "foo", Paths: []string{"a.md", "b.md"}}},
	}
	text := TextCheck(payload)
	if text == "" || !strings.Contains(text, "foo") {
		t.Fatalf("expected grouped output")
	}
}
