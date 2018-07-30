package run

import (
	"golang.org/x/sys/unix"
	"unsafe"
)

func prlimit(pid int, resource int, newlimit *unix.Rlimit,
	old *unix.Rlimit) error {
	_, _, errno := unix.RawSyscall6(unix.SYS_PRLIMIT64, uintptr(pid),
		uintptr(resource), uintptr(unsafe.Pointer(newlimit)),
		uintptr(unsafe.Pointer(old)), 0, 0)

	if errno != 0 {
		return errno
	}
	return nil
}
