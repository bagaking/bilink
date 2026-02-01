package output

import "testing"

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
