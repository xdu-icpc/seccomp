package posixtime

import "sync/atomic"

// The max number of POSIX timer.
// If try to create more timers, syscall.EAGAIN would be returned.
const LIMIT_TIMER = 8192

// This array stores channels.
var arrayChan = [LIMIT_TIMER]chan struct{}{}

// How many channels is avaliable.
var cntChan = int32(LIMIT_TIMER)

// This channel contains usable position of arrayChan.
var chanPos = make(chan int, LIMIT_TIMER)

// Atomics to mitigate ESRCH issue (#14).
var esrchFlag = [LIMIT_TIMER]int32{}

func init() {
	for i := 0; i < LIMIT_TIMER; i++ {
		chanPos <- i
	}
}

func newChanId() (int, bool) {
	value := atomic.AddInt32(&cntChan, -1)
	if value < 0 {
		atomic.AddInt32(&cntChan, 1)
		return -1, false
	}

	pos := <-chanPos
	lockUselessLock()
	arrayChan[pos] = make(chan struct{})
	esrchFlag[pos] = 0
	unlockUselessLock()
	return pos, true
}

func getChanById(id int) chan struct{} {
	lockUselessLock()
	ch := arrayChan[id]
	unlockUselessLock()
	return ch
}

func releaseChanId(id int) {
	atomic.AddInt32(&cntChan, 1)
	chanPos <- id
}
