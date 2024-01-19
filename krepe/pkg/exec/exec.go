package exec

import (
	"bytes"
	"fmt"
	"os/exec"
)

type Exec struct {
	cmd       string
	dir       string
	stdouterr *bytes.Buffer
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

func WithStdouterr(stdouterr *bytes.Buffer) Option {
	return func(e *Exec) {
		e.stdouterr = stdouterr
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
	if e.stdouterr != nil {
		cmd.Stdout = e.stdouterr
		cmd.Stderr = e.stdouterr
		err := cmd.Run()
		fmt.Println("stdouterr:", e.stdouterr.String())
		return e.stdouterr.Bytes(), err
	} else {
		return cmd.CombinedOutput()
	}
}
