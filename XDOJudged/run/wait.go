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

	wallOk := true
	if c.wallTimer != nil {
		wallOk = c.wallTimer.Stop()
	}
	cpuOk := true
	if c.cpuTimer != nil {
		cpuOk = c.cpuTimer.Stop()
	}

	// This is *critical* - from now we should *not* use any resources
	// (clock, timer, etc) bound to the PID any more.
	//
	// TODO: NeedToInvestigate: should we use the return value?
	_ = c.Cmd.Wait()

	usage := Usage{
		CPUTime: c.ProcessState.SystemTime() + c.ProcessState.UserTime(),
		MemoryInByte: 0, // FIXME: not implemented yet.
	}

	status := c.ProcessState.Sys().(syscall.WaitStatus)
	var re RuntimeError
	if status.Exited() {
		re.Code = status.ExitStatus()
	} else if status.Signaled() {
		re.Code = -int(status.Signal())
	} else {
		// TODO: what should we do?
	}

	if sigstop {
		re.Code = -int(syscall.SIGSTOP)
	}

	if !cpuOk {
		re.Reason = ReasonCPUTimeLimit
		return &usage, &re, nil
	}
	if !wallOk {
		re.Reason = ReasonWallTimeLimit
		return &usage, &re, nil
	}

	// TODO: Check for memory limit

	if re.Code == 0 {
		// No reason to believe there is a runtime error
		return &usage, nil, nil
	}
	return &usage, &re, nil
}
