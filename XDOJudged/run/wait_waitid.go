// +build linux

package run

import (
	"golang.org/x/sys/unix"
	"os"
	"runtime"
	"unsafe"
)

const (
	_P_PID       = 1
	_CLD_STOPPED = 5
)

func blockUntilWaitable(p *os.Process) (stopped bool, err error) {
	var si siginfo
	_, _, errno := unix.Syscall6(unix.SYS_WAITID, _P_PID, uintptr(p.Pid),
		uintptr(unsafe.Pointer(&si)),
		unix.WEXITED|unix.WSTOPPED|unix.WNOWAIT, 0, 0)
	runtime.KeepAlive(p)
	if errno != 0 {
		return false, errno
	}
	return si.getCode() == _CLD_STOPPED, nil
}
