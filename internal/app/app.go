package app

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/bagaking/bilink/internal/config"
	"github.com/bagaking/bilink/internal/output"
	"github.com/bagaking/bilink/internal/resolve"
	"github.com/bagaking/bilink/internal/service"
	"github.com/bagaking/bilink/internal/watch/tui"
)

func Run(args []string) error {
	if len(args) == 0 {
		return errors.New("missing command")
	}
	cmd := args[0]
	fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
	var configPath string
	var root string
	var jsonOut bool
	var interactive bool
	var noMove bool
	fs.StringVar(&configPath, "config", "", "config path")
	fs.StringVar(&root, "root", ".", "root directory")
	fs.BoolVar(&jsonOut, "json", false, "json output")
	fs.BoolVar(&interactive, "interactive", false, "interactive")
	fs.BoolVar(&noMove, "no-move", false, "do not move file")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	cfg, err := config.Load(config.ConfigOpts{Roots: []string{root}, ConfigPath: configPath})
	if err != nil {
		return err
	}
	switch cmd {
	case "refs":
		if fs.NArg() < 1 {
			return errors.New("missing target path")
		}
		payload, err := service.RunRefs(service.RefsInput{Roots: cfg.Workspace.Roots, Target: fs.Arg(0), Extensions: cfg.Scan.Extensions})
		if err != nil {
			return err
		}
		return writeOutput(jsonOut, payload, output.TextRefs(payload))
	case "check":
		payload, err := service.RunCheck(service.CheckInput{Roots: cfg.Workspace.Roots, Extensions: cfg.Scan.Extensions, ResolveRules: toResolve(cfg.Resolve), LintRules: toLintResolve(cfg.Resolve, cfg.Lint)})
		if err != nil {
			return err
		}
		return writeOutput(jsonOut, payload, output.TextCheck(payload))
	case "rename":
		if fs.NArg() < 2 {
			return errors.New("missing rename paths")
		}
		payload, err := service.RunRename(service.RenameInput{
			Roots:        cfg.Workspace.Roots,
			OldPath:      fs.Arg(0),
			NewPath:      fs.Arg(1),
			Move:         !noMove,
			Extensions:   cfg.Scan.Extensions,
			ResolveRules: toResolve(cfg.Resolve),
			Interactive:  interactive,
		})
		if err != nil {
			return err
		}
		return writeOutput(jsonOut, payload, output.TextRename(payload))
	case "watch":
		payload, err := service.RunWatch(service.WatchInput{IndexPath: cfg.Index.Path, Roots: cfg.Workspace.Roots, Extensions: cfg.Scan.Extensions})
		if err != nil {
			return err
		}
		if jsonOut {
			return writeOutput(true, payload, "")
		}
		return tui.Run(payload, configSummary(cfg))
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func writeOutput(jsonOut bool, payload any, text string) error {
	if jsonOut {
		data, err := output.JSON(payload)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}
	fmt.Println(text)
	return nil
}

func toResolve(r config.Resolve) resolve.Rules {
	return resolve.Rules{
		CaseInsensitive:      r.CaseInsensitive,
		IgnoreExtension:      r.IgnoreExtension,
		SeparatorEquivalents: r.SeparatorEquivalents,
		UnicodeNormalize:     r.UnicodeNormalize,
	}
}

func toLintResolve(base config.Resolve, lint config.Lint) resolve.Rules {
	rules := resolve.Rules{
		CaseInsensitive:      lint.RequireExactCase,
		IgnoreExtension:      lint.RequireExplicitExtension,
		SeparatorEquivalents: nil,
		UnicodeNormalize:     base.UnicodeNormalize,
	}
	if lint.RequireExactSeparators {
		rules.SeparatorEquivalents = base.SeparatorEquivalents
	}
	return rules
}

func configSummary(cfg config.Config) string {
	return fmt.Sprintf(
		"roots: %s\nextensions: %s\nresolve: caseInsensitive=%t ignoreExtension=%t separators=%v\nlint: exactCase=%t explicitExt=%t exactSeparators=%t\nindex: %s\n",
		strings.Join(cfg.Workspace.Roots, ", "),
		strings.Join(cfg.Scan.Extensions, ", "),
		cfg.Resolve.CaseInsensitive,
		cfg.Resolve.IgnoreExtension,
		cfg.Resolve.SeparatorEquivalents,
		cfg.Lint.RequireExactCase,
		cfg.Lint.RequireExplicitExtension,
		cfg.Lint.RequireExactSeparators,
		cfg.Index.Path,
	)
}
