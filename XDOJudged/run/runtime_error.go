package run

import (
	"fmt"
	"syscall"
)

// A RuntimeError represents an error occured during external program
// (compiler, object code of submitted code, checker and etc.)  Note that
// the concept is wider than an "RE" in ICPC-like contests.  Actually
// it also covers "MLE", "TLE", and "OLE".
type RuntimeError struct {
	syscall.WaitStatus
}

func (re *RuntimeError) String() string {
	if re == nil {
		return "ok, no runtime error"
	}

	if re.Exited() {
		if status := re.ExitStatus(); status != 0 {
			return fmt.Sprintf("no zero return code %d", status)
		}
	}

	if re.Signaled() {
		signal := re.Signal()
		return fmt.Sprintf("killed by signal %d (%v)", signal, signal)
	}
	return "wat the fuck?"
}
