package refs

type Link struct{ Target string }

type Index struct {
	Outbound map[string][]Link
	Inbound  map[string][]Link
}

func RefsFor(idx Index, path string) ([]Link, []Link) {
	return idx.Outbound[path], idx.Inbound[path]
}
