package siginfo

import (
	"golang.org/x/sys/unix"
	"unsafe"
)

type Siginfo struct {
	SiginfoHeader
	pad [128 - unsafe.Sizeof(SiginfoHeader{})]byte
}

type SiginfoKillHeader struct {
	SiginfoHeader
	Pid int32
	Uid uint32
}

type SiginfoKill struct {
	SiginfoKillHeader
	pad [128 - unsafe.Sizeof(SiginfoKillHeader{})]byte
}

type SiginfoTimerHeader struct {
	SiginfoHeader
	Timer, Overrun int32
	Value          uintptr
}

type SiginfoTimer struct {
	SiginfoTimerHeader
	pad [128 - unsafe.Sizeof(SiginfoTimerHeader{})]byte
}

type SiginfoQueueHeader struct {
	SiginfoHeader
	Pid   int32
	Uid   uint32
	Value uintptr
}

type SiginfoQueue struct {
	SiginfoQueueHeader
	pad [128 - unsafe.Sizeof(SiginfoQueueHeader{})]byte
}

type SiginfoChildHeader struct {
	SiginfoHeader
	Pid          int32
	Uid          uint32
	Status       unix.WaitStatus
	Utime, Stime uintptr
}

type SiginfoChild struct {
	SiginfoChildHeader
	pad [128 - unsafe.Sizeof(SiginfoChildHeader{})]byte
}

type SiginfoSyncHeader struct {
	SiginfoHeader
	Addr    uintptr
	AddrLsb int16
}

type SiginfoSync struct {
	SiginfoSyncHeader
	pad [128 - unsafe.Sizeof(SiginfoSyncHeader{})]byte
}

type SiginfoPollHeader struct {
	SiginfoHeader
	Band uintptr
	Fd   int32
}

type SiginfoPoll struct {
	SiginfoPollHeader
	pad [128 - unsafe.Sizeof(SiginfoPollHeader{})]byte
}

type SiginfoSysHeader struct {
	SiginfoHeader
	CallAddr uintptr
	Syscall  int32
	Arch     uint32
}

type SiginfoSys struct {
	SiginfoSysHeader
	pad [128 - unsafe.Sizeof(SiginfoSysHeader{})]byte
}
