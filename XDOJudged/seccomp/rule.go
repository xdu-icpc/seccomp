// Basic Seccomp rules used in XDOJ
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

import "golang.org/x/net/bpf"

// Some constants from linux/elf-em.h
const (
	EM_SPARC      = 2
	EM_386        = 3
	EM_68K        = 4
	EM_MIPS       = 8
	EM_PARISC     = 15
	EM_PPC        = 20
	EM_S390       = 22
	EM_ARM        = 40
	EM_SH         = 42
	EM_SPARCV9    = 43
	EM_IA_64      = 50
	EM_X86_64     = 62
	EM_CRIS       = 76
	EM_M32R       = 88
	EM_OPENRISC   = 92
	EM_AARCH64    = 183
	EM_TILEPRO    = 188
	EM_MICROBLAZE = 189
	EM_TILEGX     = 191
	EM_FRV        = 0x5441
	EM_ALPHA      = 0x9026
)

// Some value from linux/audit.h
const (
	AuditArch64Bit = 0x80000000
	AuditArchLE    = 0x40000000
	AuditArchN32   = 0x20000000

	AuditArchAARCH64     = EM_AARCH64 | AuditArch64Bit | AuditArchLE
	AuditArchALPHA       = EM_ALPHA | AuditArch64Bit | AuditArchLE
	AuditArchARM         = EM_ARM | AuditArchLE
	AuditArchARMEB       = EM_ARM
	AuditArchCRIS        = EM_CRIS | AuditArchLE
	AuditArchFRV         = EM_FRV
	AuditArchI386        = EM_386 | AuditArchLE
	AuditArchIA64        = EM_IA_64 | AuditArch64Bit | AuditArchLE
	AuditArchM32R        = EM_M32R
	AuditArchM68K        = EM_68K
	AuditArchMICROBLAZE  = EM_MICROBLAZE
	AuditArchMIPS        = EM_MIPS
	AuditArchMIPSEL      = EM_MIPS | AuditArchLE
	AuditArchMIPS64      = EM_MIPS | AuditArch64Bit
	AuditArchMIPS64N32   = AuditArchMIPS64 | AuditArchN32
	AuditArchMIPSEL64    = AuditArchMIPS64 | AuditArchLE
	AuditArchMIPSEL64N32 = AuditArchMIPS64N32 | AuditArchLE
	AuditArchOPENRISC    = EM_OPENRISC
	AuditArchPARISC      = EM_PARISC
	AuditArchPARISC64    = EM_PARISC | AuditArch64Bit
	AuditArchPPC         = EM_PPC
	AuditArchPPC64       = EM_PPC | AuditArch64Bit
	AuditArchPPC64LE     = AuditArchPPC64 | AuditArchLE
	AuditArchS390        = EM_S390
	AuditArchS390X       = EM_S390 | AuditArch64Bit
	AuditArchSH          = EM_SH
	AuditArchSHEL        = EM_SH | AuditArchLE
	AuditArchSH64        = EM_SH | AuditArch64Bit
	AuditArchSHEL64      = AuditArchSHEL | AuditArch64Bit
	AuditArchSPARC       = EM_SPARC
	AuditArchSPARC64     = EM_SPARCV9 | AuditArch64Bit
	AuditArchTILEGX      = EM_TILEGX | AuditArch64Bit | AuditArchLE
	AuditArchTILEGX32    = EM_TILEGX | AuditArchLE
	AuditArchTILEPRO     = EM_TILEPRO | AuditArchLE
	AuditArchX86_64      = EM_X86_64 | AuditArch64Bit | AuditArchLE
)

// Some value from linux/seccomp.h
const (
	SECCOMP_RET_KILL_PROCESS = 0x80000000
	SECCOMP_RET_KILL_THREAD  = 0x00000000
	SECCOMP_RET_KILL         = SECCOMP_RET_KILL_THREAD
	SECCOMP_RET_TRAP         = 0x00030000
	SECCOMP_RET_ERRNO        = 0x00050000
	SECCOMP_RET_TRACE        = 0x7ff00000
	SECCOMP_RET_LOG          = 0x7ffc0000
	SECCOMP_RET_ALLOW        = 0x7fff0000
)

// Generate an instruction loading AUDIT_ARCH_* value from seccomp data.
func LoadArch() bpf.Instruction {
	return bpf.LoadAbsolute{Off: 4, Size: 4}
}

// Generate an instruction loading the system call number from seccomp data.
func LoadNr() bpf.Instruction {
	return bpf.LoadAbsolute{Off: 0, Size: 4}
}

// Mark a half of a 64-bit register.
type Half uint32

const (
	Low  Half = 0
	High Half = 4
)

// Generate an instruction to load one half of instruction pointer from
// seccomp data.
func LoadIP(h Half) bpf.Instruction {
	return bpf.LoadAbsolute{Off: 8 + uint32(h), Size: 4}
}

// Generate an instruction to load one half of a 64-bit register from
// seccomp data.
func LoadReg(n uint32, h Half) bpf.Instruction {
	return bpf.LoadAbsolute{Off: 8 + n*8 + uint32(h), Size: 4}
}

// Generate an instruction to allow the system call.
func RetAllow() bpf.Instruction {
	return bpf.RetConstant{Val: SECCOMP_RET_ALLOW}
}

// Generate an instruction to kill the process.
func RetKillProcess() bpf.Instruction {
	return bpf.RetConstant{Val: SECCOMP_RET_KILL_PROCESS}
}

// Generate an instruction to kill the thread.
func RetKillThread() bpf.Instruction {
	return bpf.RetConstant{Val: SECCOMP_RET_KILL_THREAD}
}

// Generate an instruction to send a catchable SIGSYS signal to the process,
// and set the si_errno field of the siginfo_t structure to errno.
func RetTrap(errno uint16) bpf.Instruction {
	return bpf.RetConstant{Val: SECCOMP_RET_TRAP | uint32(errno)}
}

// Generate an instruction to return an errno for the system call, without
// executing it.
func RetErrno(errno uint16) bpf.Instruction {
	return bpf.RetConstant{Val: SECCOMP_RET_ERRNO | uint32(errno)}
}

// Generate an instruction to notify a ptrace-based tracer prior to
// executing the system call.  The tracer can use PTRACE_GETEVENTMSG
// to get the msg value.
func RetTrace(msg uint16) bpf.Instruction {
	return bpf.RetConstant{Val: SECCOMP_RET_TRACE | uint32(msg)}
}

// Allow the system call being executed after the filter return action
// is logged.
func RetLog() bpf.Instruction {
	return bpf.RetConstant{Val: SECCOMP_RET_LOG}
}
