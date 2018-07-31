// XDOJ Seccomp rule for x86 and x86-64.
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

// +build amd64 386

package seccomp

import (
	. "golang.org/x/net/bpf"
	"golang.org/x/sys/unix"
)

const x32SyscallBit = uint32(0x40000000)

var noForkRule = []Instruction{
	// 0, architecture check
	LoadArch(),
	// 1, if not X86_64 skip to I386 check
	JumpIf{Cond: JumpEqual,
		Val:       AuditArchX86_64,
		SkipTrue:  0,
		SkipFalse: 5},
	// 2, load syscall number
	LoadNr(),
	// 3, remove the x32 syscall bit
	ALUOpConstant{Op: ALUOpAnd, Val: ^x32SyscallBit},
	// 4, if syscall is fork, return errno
	JumpIf{Cond: JumpEqual,
		Val:       57,
		SkipTrue:  9,
		SkipFalse: 0},
	// 5, if syscall is vfork, return errno
	JumpIf{Cond: JumpEqual,
		Val:       58,
		SkipTrue:  8,
		SkipFalse: 0},
	// 6, if syscall is clone, go to CLONE_THREAD flag check, otherwise
	// return OK
	JumpIf{Cond: JumpEqual,
		Val:       56,
		SkipTrue:  4,
		SkipFalse: 6},
	// 7, load syscall number for I386
	LoadNr(),
	// 8, if syscall is fork, return errno
	JumpIf{Cond: JumpEqual,
		Val:       2,
		SkipTrue:  5,
		SkipFalse: 0},
	// 9, if syscall is vfork, return errno
	JumpIf{Cond: JumpEqual,
		Val:       190,
		SkipTrue:  4,
		SkipFalse: 0},
	// 10, if syscall is clone, do CLONE_THREAD flag check, otherwise
	// return OK
	JumpIf{Cond: JumpEqual,
		Val:       120,
		SkipTrue:  0,
		SkipFalse: 2},
	// 11, CLONE_THREAD flag check, load the first argument, see clone(2)
	// "C library/kernel differences".
	LoadReg(1, Low),
	// 12, if CLONE_THREAD is set, return OK, otherwise return errno
	JumpIf{Cond: JumpBitsSet,
		Val:       unix.CLONE_THREAD,
		SkipTrue:  0,
		SkipFalse: 1},
	// 13, return OK
	RetOK,
	// 14, return ENOSYS
	RetDisallow,
}
