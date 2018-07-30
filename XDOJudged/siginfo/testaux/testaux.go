package testaux

/*
#include <signal.h>

const int offset_pid = __builtin_offsetof(siginfo_t, si_pid);
const int offset_uid = __builtin_offsetof(siginfo_t, si_uid);
const int offset_timer = __builtin_offsetof(siginfo_t, si_timerid);
const int offset_overrun = __builtin_offsetof(siginfo_t, si_overrun);
const int offset_status = __builtin_offsetof(siginfo_t, si_status);
const int offset_value = __builtin_offsetof(siginfo_t, si_value);
const int offset_addr = __builtin_offsetof(siginfo_t, si_addr);
const int offset_addr_lsb = __builtin_offsetof(siginfo_t, si_addr_lsb);
const int offset_band = __builtin_offsetof(siginfo_t, si_band);
const int offset_fd = __builtin_offsetof(siginfo_t, si_fd);
const int offset_call_addr = __builtin_offsetof(siginfo_t, si_call_addr);
const int offset_syscall = __builtin_offsetof(siginfo_t, si_syscall);
const int offset_arch = __builtin_offsetof(siginfo_t, si_arch);
const int offset_utime = __builtin_offsetof(siginfo_t, si_utime);
const int offset_stime = __builtin_offsetof(siginfo_t, si_stime);
*/
import "C"

const (
	OffsetPid = C.offset_pid
	OffsetUid = C.offset_uid
	OffsetTimer = C.offset_timer
	OffsetOverrun = C.offset_overrun
	OffsetStatus = C.offset_status
	OffsetValue = C.offset_value
	OffsetAddr = C.offset_addr
	OffsetAddrLsb = C.offset_addr_lsb
	OffsetBand = C.offset_band
	OffsetFd = C.offset_fd
	OffsetCallAddr = C.offset_call_addr
	OffsetSyscall = C.offset_syscall
	OffsetArch = C.offset_arch
	OffsetUtime = C.offset_utime
	OffsetStime = C.offset_stime
)
