package exec

import (
	"fmt"
	"os/exec"
)

type Exec struct {
	cmd string
	dir string
}

type Option func(*Exec)

func WithCmd(cmd string) Option {
	return func(e *Exec) {
		e.cmd = cmd
	}
}

func WithDir(dir string) Option {
	return func(e *Exec) {
		e.dir = dir
	}
}

func NewExec(options ...Option) *Exec {
	exec := &Exec{}
	for _, option := range options {
		option(exec)
	}
	return exec
}

func (e *Exec) Run(args ...string) ([]byte, error) {
	cmd := exec.Command(e.cmd, args...)
	cmd.Dir = e.dir
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	return out, err
}
