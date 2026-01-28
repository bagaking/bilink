package check

type File struct {
	Path       string
	ResolveKey string
	LintKey    string
}

func Detect(files []File) ([]string, []string) {
	seen := map[string][]string{}
	for _, f := range files {
		seen[f.ResolveKey] = append(seen[f.ResolveKey], f.Path)
	}
	var errs []string
	for key, paths := range seen {
		if len(paths) > 1 {
			errs = append(errs, key)
		}
	}
	return errs, nil
}
