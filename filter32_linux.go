// Copyright (C) 2019  Laboratory of ICPC, Xidian University
// SPDX-License-Identifier: AGPL-3.0-or-later

// Author: Xi Ruoyao <xry111@mengyan1223.wang>

// +build arm 386 mips mipsle

package seccomp

import "golang.org/x/net/bpf"

type SockFprog struct {
	Len       uint16
	pad_cgo_0 [2]byte
	Filter    *bpf.RawInstruction
}
