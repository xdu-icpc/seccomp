// Unit test for SeccompFilter.
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

package seccomp_test

import (
	"golang.org/x/net/bpf"
	"golang.org/x/sys/unix"
	"runtime"
	"testing"
	"unsafe"

	"linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/seccomp"
)

func TestSeccompFilter(t *testing.T) {
	// We use SYS_GETCPU for test purpose.  Hope Go runtime doesn't
	// need it.
	rule := []bpf.Instruction{
		seccomp.LoadNr(),
		bpf.JumpIf{
			Cond:      bpf.JumpEqual,
			Val:       unix.SYS_GETCPU,
			SkipTrue:  0,
			SkipFalse: 1,
		},
		seccomp.RetErrno(uint16(unix.ENOSYS)),
		seccomp.RetAllow(),
	}
	filter, err := bpf.Assemble(rule)
	if err != nil {
		t.Fatalf("can not assemble bpf rule: %v", err)
	}

	// Filter is per-thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err = seccomp.SeccompFilter(0, filter)
	if err != nil {
		t.Fatal(err)
	}

	var cpu, node uint32
	_, _, err = unix.RawSyscall(unix.SYS_GETCPU,
		uintptr(unsafe.Pointer(&cpu)), uintptr(unsafe.Pointer(&node)), 0)
	if err != unix.ENOSYS {
		t.Fatalf("GETCPU didn't fail with ENOSYS: %v", err)
	}
}
