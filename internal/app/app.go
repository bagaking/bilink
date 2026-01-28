package app

import "fmt"

func Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing command")
	}
	switch args[0] {
	case "refs", "rename", "check", "watch":
		return nil
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}
