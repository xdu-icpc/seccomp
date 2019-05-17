package seccomp

import (
	"golang.org/x/net/bpf"
	"syscall"
)

// Create a SockFprog instance which can be used with seccomp syscall.
func NewSockFprog(insn []bpf.RawInstruction) (*SockFprog, error) {
	if insn == nil {
		return nil, nil
	}

	l := len(insn)
	if l > 4096 {
		return nil, syscall.EINVAL
	}

	return &SockFprog{
		Len: uint16(l),
		Filter: &insn[0],
	}, nil
}
