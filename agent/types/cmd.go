package types

type CommandExec func(program string, args ...string) error
