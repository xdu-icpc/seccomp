package siginfo

import (
	"testing"
	"unsafe"

	. "linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/siginfo/testaux"
)

func assertEqual(t *testing.T, expect uintptr, get uintptr) {
	t.Helper()
	if expect != get {
		t.Errorf("result mismatch: expect %v, get %v", expect, get)
	}
}

func TestSize(t *testing.T) {
	assertEqual(t, unsafe.Sizeof(SiginfoKill{}), 128)
	assertEqual(t, unsafe.Sizeof(SiginfoTimer{}), 128)
	assertEqual(t, unsafe.Sizeof(SiginfoQueue{}), 128)
	assertEqual(t, unsafe.Sizeof(SiginfoSync{}), 128)
	assertEqual(t, unsafe.Sizeof(SiginfoPoll{}), 128)
	assertEqual(t, unsafe.Sizeof(SiginfoSys{}), 128)
}

func TestABI(t *testing.T) {
	assertEqual(t, unsafe.Offsetof(SiginfoKill{}.Pid), OffsetPid)
	assertEqual(t, unsafe.Offsetof(SiginfoKill{}.Uid), OffsetUid)
	assertEqual(t, unsafe.Offsetof(SiginfoChild{}.Status), OffsetStatus)
	assertEqual(t, unsafe.Offsetof(SiginfoTimer{}.Timer), OffsetTimer)
	assertEqual(t, unsafe.Offsetof(SiginfoTimer{}.Value), OffsetValue)
	assertEqual(t, unsafe.Offsetof(SiginfoTimer{}.Overrun), OffsetOverrun)
	assertEqual(t, unsafe.Offsetof(SiginfoChild{}.Utime), OffsetUtime)
	assertEqual(t, unsafe.Offsetof(SiginfoChild{}.Stime), OffsetStime)
	t.Logf("SI_TIMER = %d", SI_TIMER)
}
