// Seccomp syscall.
// Copyright (C) 2017-2018  Laboratory of ACM/ICPC, Xidian University

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

package seccomp

import (
	"golang.org/x/net/bpf"
	"golang.org/x/sys/unix"
	"unsafe"
)

const (
	SECCOMP_FILTER_FLAG_TSYNC = 1
	SECCOMP_FILTER_FLAG_LOG   = 2
)

// I don't think "strict mode" can support Go program.

// Install a seccomp filter.
func SeccompFilter(flags uintptr, filter []bpf.RawInstruction) error {
	if len(filter) > 255 {
		return unix.EINVAL
	}

	sockProg := sockFprog{
		Len:    uint16(len(filter)),
		Filter: &filter[0],
	}

	_, _, errno := unix.RawSyscall(unix.SYS_SECCOMP, 1, flags,
		uintptr(unsafe.Pointer(&sockProg)))
	if errno != 0 {
		return errno
	}

	return nil
}
