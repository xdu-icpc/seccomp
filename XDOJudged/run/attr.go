package run

import (
	"golang.org/x/sys/unix"
	"time"

	"linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/bind"
)

type ResourceLimit struct {
	Resource int
	Rlimit   unix.Rlimit
}

// Attr holds the attributes that will be applied to a new process started
// by this package.
type Attr struct {
	// If this field is not nil, a seccomp filter will be used to prevent
	// forking since we have no way to limit a process group's CPU time.
	CPUTimeLimit *time.Duration

	WallTimeLimit *time.Duration

	ResourceLimit []ResourceLimit

	BindMount []bind.BindMount

	// TODO: Fields for cgroups...
}
