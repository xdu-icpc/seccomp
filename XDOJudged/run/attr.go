package run

import (
	"golang.org/x/sys/unix"
	"strings"
	"time"
)

type ResourceLimit struct {
	Resource int
	Rlimit   unix.Rlimit
}

type BindMount struct {
	OldDir, NewDir string
	NoRecursive, ReadOnly bool
}

func (bind *BindMount) String() string {
	path := []string{bind.OldDir, bind.NewDir, "noro", "rbind"}
	if bind.ReadOnly {
		path[2] = "ro"
	}
	if bind.NoRecursive {
		path[3] = "norbind"
	}
	return strings.Join(path, ":")
}


// Attr holds the attributes that will be applied to a new process started
// by this package.
type Attr struct {
	// If this field is not nil, a seccomp filter will be used to prevent
	// forking since we have no way to limit a process group's CPU time.
	CPUTimeLimit *time.Duration

	WallTimeLimit *time.Duration

	ResourceLimit []ResourceLimit

	BindMount []BindMount

	// TODO: Fields for cgroups...
}
