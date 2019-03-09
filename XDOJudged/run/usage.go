package run

import (
	"fmt"
	"time"
)

// A Usage contains runtime usage info of a process.
type Usage struct {
	// CPU time consumed by the process.
	CPUTime time.Duration
	// Memory consumed by the process.  A zero value means we failed
	// to get it.
	MemoryInByte int64
}

func (u *Usage) String() string {
	return fmt.Sprintf("CPU time = %v, memory = %v byte(s)",
		u.CPUTime, u.MemoryInByte)
}
