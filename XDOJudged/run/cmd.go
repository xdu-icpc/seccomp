// High level package for XDOJ sandbox function.
package run

import (
	"os/exec"
	"time"

	"linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/posixtime"
)

// Cmd represents an external command being prepared or run.
//
// A Cmd cannot be reused after calling its Start, Run, Output or
// CombinedOutput methods.
type Cmd struct {
	*exec.Cmd
	Attr *Attr

	cpuTimer  *posixtime.Timer
	wallTimer *time.Timer
	finished  bool
}

// Run starts the specified command and waits for it to complete.
func (c *Cmd) Run() (*Usage, *RuntimeError, error) {
	err := c.Start()
	if err != nil {
		return nil, nil, err
	}
	return c.Wait()
}

// Start starts the specified command but does not wait for it to complete.
func (c *Cmd) Start() error {
	return c.start()
}

// Wait waits for the command to exit and waits for any copying to stdin or
// copying from stdout or stderr to complete.
func (c *Cmd) Wait() (*Usage, *RuntimeError, error) {
	return c.wait()
}

// Command is a replica of (os/exec).Command.
func Command(name string, arg ...string) *Cmd {
	return &Cmd{Cmd: exec.Command(name, arg...)}
}
