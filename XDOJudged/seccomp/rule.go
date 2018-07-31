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

import (
	. "debug/elf"
	"golang.org/x/net/bpf"
	"golang.org/x/sys/unix"
)

// Some value from linux/audit.h
const (
	auditArch64Bit = 0x80000000
	auditArchLE    = 0x40000000
	auditArchN32   = 0x20000000

	AuditArchAARCH64     = uint32(EM_AARCH64) | auditArch64Bit | auditArchLE
	AuditArchALPHA       = uint32(EM_ALPHA) | auditArch64Bit | auditArchLE
	AuditArchARM         = uint32(EM_ARM) | auditArchLE
	AuditArchARMEB       = uint32(EM_ARM)
	AuditArchCRIS        = uint32(76) | auditArchLE
	AuditArchFRV         = uint32(0x5441)
	AuditArchI386        = uint32(EM_386) | auditArchLE
	AuditArchIA64        = uint32(EM_IA_64) | auditArch64Bit | auditArchLE
	AuditArchM32R        = uint32(88)
	AuditArchM68K        = uint32(EM_68K)
	AuditArchMICROBLAZE  = uint32(189)
	AuditArchMIPS        = uint32(EM_MIPS)
	AuditArchMIPSEL      = uint32(EM_MIPS) | auditArchLE
	AuditArchMIPS64      = uint32(EM_MIPS) | auditArch64Bit
	AuditArchMIPS64N32   = AuditArchMIPS64 | auditArchN32
	AuditArchMIPSEL64    = AuditArchMIPS64 | auditArchLE
	AuditArchMIPSEL64N32 = AuditArchMIPS64N32 | auditArchLE
	AuditArchOPENRISC    = uint32(92)
	AuditArchPARISC      = uint32(EM_PARISC)
	AuditArchPARISC64    = uint32(EM_PARISC) | auditArch64Bit
	AuditArchPPC         = uint32(EM_PPC)
	AuditArchPPC64       = uint32(EM_PPC) | auditArch64Bit
	AuditArchPPC64LE     = AuditArchPPC64 | auditArchLE
	AuditArchS390        = uint32(EM_S390)
	AuditArchS390X       = uint32(EM_S390) | auditArch64Bit
	AuditArchSH          = uint32(EM_SH)
	AuditArchSHEL        = uint32(EM_SH) | auditArchLE
	AuditArchSH64        = uint32(EM_SH) | auditArch64Bit
	AuditArchSHEL64      = AuditArchSHEL | auditArch64Bit
	AuditArchSPARC       = uint32(EM_SPARC)
	AuditArchSPARC64     = uint32(EM_SPARCV9) | auditArch64Bit
	AuditArchTILEGX      = uint32(191) | auditArch64Bit | auditArchLE
	AuditArchTILEGX32    = uint32(191) | auditArchLE
	AuditArchTILEPRO     = uint32(188) | auditArchLE
	AuditArchX86_64      = uint32(EM_X86_64) | auditArch64Bit | auditArchLE
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

// 32-bit field load instructions
var (
	LoadArch = bpf.LoadAbsolute{Off: 4, Size: 4}
	LoadNr   = bpf.LoadAbsolute{Off: 0, Size: 4}
)

// 64-bit field load instructions.
// For Little Endian.  BE should revert H(igh)/L(ow).
var (
	LoadIPLow  = bpf.LoadAbsolute{Off: 8, Size: 4}
	LoadIPHigh = bpf.LoadAbsolute{Off: 12, Size: 4}
	LoadA1Low  = bpf.LoadAbsolute{Off: 16, Size: 4}
	LoadA1High = bpf.LoadAbsolute{Off: 20, Size: 4}
	LoadA2Low  = bpf.LoadAbsolute{Off: 24, Size: 4}
	LoadA2High = bpf.LoadAbsolute{Off: 28, Size: 4}
	LoadA3Low  = bpf.LoadAbsolute{Off: 32, Size: 4}
	LoadA3High = bpf.LoadAbsolute{Off: 36, Size: 4}
	LoadA4Low  = bpf.LoadAbsolute{Off: 40, Size: 4}
	LoadA4High = bpf.LoadAbsolute{Off: 44, Size: 4}
	LoadA5Low  = bpf.LoadAbsolute{Off: 48, Size: 4}
	LoadA5High = bpf.LoadAbsolute{Off: 52, Size: 4}
	LoadA6Low  = bpf.LoadAbsolute{Off: 56, Size: 4}
	LoadA6High = bpf.LoadAbsolute{Off: 60, Size: 4}
)

// The abbr. for CLONE_THREAD
const tflag = uint32(unix.CLONE_THREAD)

// Actions
var (
	RetOK       = bpf.RetConstant{Val: SECCOMP_RET_ALLOW}
	RetDisallow = bpf.RetConstant{Val: SECCOMP_RET_KILL}
)

// A seccomp prohibits process creating, including fork(2), vfork(2)
// and clone(2) without CLONE_THREAD flag.
var NoForkFilter []bpf.RawInstruction

func init() {
	var err error
	NoForkFilter, err = bpf.Assemble(noForkRule)
	if err != nil {
		panic(err)
	}
}
