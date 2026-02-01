package config

import (
	"errors"
	"os"

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
	if opts.ConfigPath == "" {
		return cfg, nil
	}
	data, err := os.ReadFile(opts.ConfigPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return Config{}, err
	}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
