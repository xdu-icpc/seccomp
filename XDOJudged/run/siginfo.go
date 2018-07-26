package run

/*
#include <signal.h>

const int offset_of_si_code = __builtin_offsetof(siginfo_t, si_code);
*/
import "C"

import "unsafe"

type siginfo C.siginfo_t

func (i *siginfo) getCode() int32 {
	// Some platforms have si_code in anonymous union, so we have to use
	// a dirty expression.
	return *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(i)) +
		uintptr(C.offset_of_si_code)))
}
