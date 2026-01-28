package index

import "testing"

func TestStoreRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/index.json"
	idx := Index{Outbound: map[string][]Link{"a.md": {{Target: "b"}}}}
	if err := Save(path, idx); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Outbound["a.md"]) != 1 {
		t.Fatalf("unexpected index")
	}
}
