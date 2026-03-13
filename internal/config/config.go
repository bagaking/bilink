package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Workspace Workspace `toml:"workspace"`
	Scan      Scan      `toml:"scan"`
	Resolve   Resolve   `toml:"resolveRules"`
	Lint      Lint      `toml:"lintRules"`
	Anchors   Anchors   `toml:"anchors"`
	Update    Update    `toml:"updatePolicy"`
	Index     Index     `toml:"index"`
}

type Workspace struct {
	Roots []string `toml:"roots"`
}

type Scan struct {
	Extensions []string `toml:"extensions"`
}

type Resolve struct {
	CaseInsensitive      bool     `toml:"caseInsensitive"`
	IgnoreExtension      bool     `toml:"ignoreExtension"`
	SeparatorEquivalents []string `toml:"separatorEquivalents"`
	UnicodeNormalize     string   `toml:"unicodeNormalize"`
}

type Lint struct {
	RequireExactCase         bool `toml:"requireExactCase"`
	RequireExplicitExtension bool `toml:"requireExplicitExtension"`
	RequireExactSeparators   bool `toml:"requireExactSeparators"`
}

type Anchors struct {
	Style string `toml:"style"`
	Mode  string `toml:"mode"`
}

type Update struct {
	Mode string `toml:"mode"`
}

type Index struct {
	Path             string `toml:"path"`
	RequiredForWatch bool   `toml:"requiredForWatch"`
}

type ConfigOpts struct {
	Roots      []string
	ConfigPath string
}

func Defaults() Config {
	return Config{
		Workspace: Workspace{Roots: []string{"."}},
		Scan:      Scan{Extensions: []string{".md", ".markdown", ".mdx"}},
		Resolve: Resolve{
			CaseInsensitive:      true,
			IgnoreExtension:      true,
			SeparatorEquivalents: []string{" ", "-", "_"},
			UnicodeNormalize:     "NFKC",
		},
		Lint: Lint{
			RequireExactCase:         true,
			RequireExplicitExtension: true,
			RequireExactSeparators:   true,
		},
		Anchors: Anchors{Style: "github", Mode: "resolve-only"},
		Update:  Update{Mode: "balanced"},
		Index:   Index{Path: ".bilink/index.json", RequiredForWatch: true},
	}
}

func Load(opts ConfigOpts) (Config, error) {
	cfg := Defaults()
	if len(opts.Roots) > 0 {
		cfg.Workspace.Roots = opts.Roots
	}
	configPath := opts.ConfigPath
	explicitConfig := configPath != ""
	configRoot := ""
	if configPath == "" {
		root := "."
		if len(cfg.Workspace.Roots) > 0 {
			root = cfg.Workspace.Roots[0]
		}
		configRoot = root
		configPath = filepath.Join(root, ".bilink", "settings.toml")
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) && !explicitConfig {
			return cfg, nil
		}
		return Config{}, err
	}
	var fileRoots struct {
		Workspace struct {
			Roots []string `toml:"roots"`
		} `toml:"workspace"`
	}
	if err := toml.Unmarshal(data, &fileRoots); err != nil {
		return Config{}, err
	}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	if !explicitConfig && configRoot != "" && fileRoots.Workspace.Roots != nil {
		cfg.Workspace.Roots = resolveDefaultConfigRoots(configRoot, cfg.Workspace.Roots)
	}
	return cfg, nil
}

func resolveDefaultConfigRoots(configRoot string, roots []string) []string {
	resolved := make([]string, 0, len(roots))
	for _, root := range roots {
		if filepath.IsAbs(root) {
			resolved = append(resolved, root)
			continue
		}
		resolved = append(resolved, filepath.Join(configRoot, root))
	}
	return resolved
}
