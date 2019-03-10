package run

import "syscall"

func (c *Cmd) wait() (*Usage, *RuntimeError, error) {
	if c.Process == nil || c.finished {
		return nil, nil, c.Cmd.Wait()
	}
	c.finished = true

	sigstop := false
	for {
		stopped, err := blockUntilWaitable(c.Process)
		if err != nil {
			return nil, nil, err
		}
		if stopped {
			sigstop = true
		} else {
			break
		}
		c.kill()
	}

	// This is *critical* - from now we should *not* use any resources
	// (clock, timer, etc) bound to the PID any more.
	//
	// TODO: NeedToInvestigate: should we use the return value?
	_ = c.Cmd.Wait()

	rusage, ok := c.ProcessState.SysUsage().(*syscall.Rusage)
	memusage := int64(0)
	if ok {
		memusage = rusage.Maxrss * 1024
	}

	usage := Usage{
		CPUTime:      c.ProcessState.SystemTime() + c.ProcessState.UserTime(),
		MemoryInByte: memusage,
	}

	status := c.ProcessState.Sys().(syscall.WaitStatus)
	re := &RuntimeError{WaitStatus: status}

	if sigstop {
		// TODO
	}

	if re.Exited() && re.ExitStatus() == 0 {
		// No reason to believe there is a runtime error
		return &usage, nil, nil
	}
	return &usage, re, nil
}
