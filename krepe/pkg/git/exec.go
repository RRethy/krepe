package git

import (
	"os/exec"
)

type Exec struct {
	dir string
}

func NewExec(dir string) *Exec {
	return &Exec{dir: dir}
}

func (e *Exec) Run(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = e.dir
	return cmd.CombinedOutput()
}
