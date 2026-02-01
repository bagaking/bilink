package service

type CommandKind int

const (
	CommandRefs CommandKind = iota
	CommandCheck
	CommandRename
	CommandWatch
)

func (c CommandKind) String() string {
	switch c {
	case CommandRefs:
		return "refs"
	case CommandCheck:
		return "check"
	case CommandRename:
		return "rename"
	case CommandWatch:
		return "watch"
	default:
		return "unknown"
	}
}

type Options struct {
	Roots       []string
	ConfigPath  string
	JSON        bool
	Interactive bool
	Move        bool
	OldPath     string
	NewPath     string
	TargetPath  string
}
