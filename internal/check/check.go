package check

type File struct {
	Path       string
	ResolveKey string
	LintKey    string
}

func Detect(files []File) ([]string, []string) {
	resolveSeen := map[string][]string{}
	lintSeen := map[string][]string{}
	for _, f := range files {
		resolveSeen[f.ResolveKey] = append(resolveSeen[f.ResolveKey], f.Path)
		if f.LintKey != "" {
			lintSeen[f.LintKey] = append(lintSeen[f.LintKey], f.Path)
		}
	}
	var errs []string
	for key, paths := range resolveSeen {
		if len(paths) > 1 {
			errs = append(errs, key)
		}
	}
	var warns []string
	for key, paths := range lintSeen {
		if len(paths) > 1 {
			warns = append(warns, key)
		}
	}
	return errs, warns
}
