package run

import (
	"fmt"
	"time"
)

// A Usage contains runtime usage info of a process.
type Usage struct{
	CPUTime time.Duration
	MemoryInByte int64
}

func (u *Usage) String() string {
	return fmt.Sprintf("CPU time = %v, memory = %v byte(s)",
		u.CPUTime, u.MemoryInByte)
}
