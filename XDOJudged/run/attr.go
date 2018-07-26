package run

import (
	"time"
)

// Attr holds the attributes that will be applied to a new process started
// by this package.
type Attr struct{
	CPUTimeLimit *time.Duration
	WallTimeLimit *time.Duration
	// TODO: Fields for cgroups...
}
