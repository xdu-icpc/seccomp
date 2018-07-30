package siginfo

import (
	"syscall"
	"unsafe"
)

func (s *SiginfoHeader) Signal() syscall.Signal {
	return syscall.Signal(s.Signo)
}

func (s *Siginfo) ToSiginfoChild() *SiginfoChild {
	if s.Signal() == syscall.SIGCHLD {
		return (*SiginfoChild)(unsafe.Pointer(s))
	}
	return nil
}

func (s *Siginfo) ToSiginfoTimer() *SiginfoTimer {
	if s.Code == SI_TIMER {
		return (*SiginfoTimer)(unsafe.Pointer(s))
	}
	return nil
}

func (s *Siginfo) ToSiginfoQueue() *SiginfoQueue {
	if s.Code == SI_QUEUE {
		return (*SiginfoQueue)(unsafe.Pointer(s))
	}
	return nil
}

func (s *Siginfo) ToSiginfoPoll() *SiginfoPoll {
	if s.Signal() == syscall.SIGPOLL {
		return (*SiginfoPoll)(unsafe.Pointer(s))
	}
	return nil
}

func (s *Siginfo) ToSiginfoSys() *SiginfoSys {
	if s.Signal() == syscall.SIGSYS {
		return (*SiginfoSys)(unsafe.Pointer(s))
	}
	return nil
}
