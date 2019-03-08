package bind

import (
	"fmt"
	"strings"
)

type BindMount struct {
	OldDir, NewDir string
	NoRecursive, ReadOnly bool
}

func (b *BindMount) String() string {
	path := []string{b.OldDir, b.NewDir, "noro", "rbind"}
	if b.ReadOnly {
		path[2] = "ro"
	}
	if b.NoRecursive {
		path[3] = "norbind"
	}
	return strings.Join(path, ":")
}

func Parse(s string) (*BindMount, error) {
	path := strings.Split(s, ":")
	// "<old>:<new>:<ro>:<rbind>"
	if len(path) != 4 {
		return nil, fmt.Errorf("can not parse %s", s)
	}
	return &BindMount{
		OldDir: path[0],
		NewDir: path[1],
		ReadOnly: path[2] == "ro",
		NoRecursive: path[3] != "rbind",
	}, nil;
}
