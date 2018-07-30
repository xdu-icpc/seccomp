// Wrap signal related syscalls of Linux.
// Copyright (C) 2017  Laboratory of ACM/ICPC, Xidian University

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Author: Xi Ruoyao <ryxi@stu.xidian.edu.cn>

// +build linux

package posixtime

import (
	"golang.org/x/sys/unix"
	"os"
	"syscall"
	"unsafe"

	"linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/siginfo"
)

// In Linux, sigwaitinfo is just a special case of rt_sigtimedwait,
// with sigsetsize = sizeof(sigset_t) and timeout = NULL. To reduce GC we
// do not return siginfo, but accept a pointer as parameter.
func sigwaitinfo(set *sigset, info *siginfo.Siginfo) error {
	_, _, errno := unix.Syscall6(unix.SYS_RT_SIGTIMEDWAIT,
		uintptr(unsafe.Pointer(set)), uintptr(unsafe.Pointer(info)),
		uintptr(0), uintptr(8), uintptr(0), uintptr(0))

	if errno != 0 {
		return errno
	}

	return nil
}

func sigqueue(pid int, sig syscall.Signal, value uintptr) error {
	si := &siginfo.SiginfoQueue{}
	si.Signo = int32(sig)
	si.Code = int32(siginfo.SI_QUEUE)
	si.Pid = int32(os.Getpid())
	si.Uid = uint32(os.Getuid())
	si.Value = value
	_, _, errno := unix.Syscall(uintptr(pid), uintptr(sig),
		uintptr(unsafe.Pointer(si)), 0)

	if errno != 0 {
		return errno
	}

	return nil
}
