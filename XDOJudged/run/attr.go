package run

import (
	"golang.org/x/sys/unix"

	"linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/bind"
)

// Attr holds the attributes that will be applied to a new process started
// by this package.
type Attr struct {
	// If this is true, we won't drop (real or namespace) root's
	// capabilities.
	KeepCap bool

	// Bind mount filesystems.  If it is not empty but CLONE_NEWNS is not
	// set, we'll refuse to modify current namespace and return an error.
	BindMount []bind.BindMount

	// Seccomp filter to use (in JSON encoded []bpf.Instruction)
	Seccomp string
}
