// Low level routine implementing key timer functions.
// Copyright (C) 2017  Laboratory of ACM/ICPC, Xidian University

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Author: Xi Ruoyao <ryxi@stu.xidian.edu.cn>

package posixtime

import (
	"os"
	"sync/atomic"
	"syscall"

	"golang.org/x/sys/unix"
)

var zeroits = itimerspec{}

func timerRoutine(c ClockID, ts unix.Timespec, chstop chan struct{},
	callback func()) (err error) {
	chexpireId, ok := newChanId()
	if !ok { // too many timers
		return syscall.EAGAIN
	}
	chexpire := getChanById(chexpireId)
	defer func() {
		if err != nil {
			close(chexpire)
			releaseChanId(chexpireId)
		}
	}()

	// create the timer
	sigev := sigevent{
		sigev_notify: _SIGEV_SIGNAL,
		sigev_signo:  _SIGRTMIN,
	}
	sigev.setValue(uintptr(chexpireId))

	t, err := timerCreate(c, &sigev)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			timerDelete(t)
		}
	}()

	// set time for the timer
	its := itimerspec{
		interval: unix.Timespec{},
		value:    ts,
	}
	_, err = timerSetTime(t, 0, &its)
	if err != nil {
		return err
	}

	go func() {
		select {
		case <-chstop:
			// stop the timer by setting a zero itimerspec.
			its, err := timerSetTime(t, 0, &zeroits)
			if err != nil {
				if err == unix.ESRCH {
					// The process bound to the timer has been reaped.
					old := atomic.SwapInt32(&esrchFlag[chexpireId], 1)
					if old == 0 {
						sigqueue(os.Getpid(), SIGRTMIN,
							uintptr(chexpireId))
					}
				} else {
					panic(err)
				}
			}

			if its == nil || *its == zeroits {
				// This means the timer has already expired.
				// There is a potential race conditon: the expiration
				// signal is still in queue and chexpire is NOT closed
				// yet.  Then we can't call releaseChanId, or demux()
				// may send to a channel without receiver and block forever.
				select {
				case <-chexpire: // Now we know chexpire is closed.
				}
			} else {
				// Otherwise, the timer has been stopped successfully
				// before it could expire.  Though we can let GC close it,
				// we do this ourselves to make this explicitly.
				close(chexpire)
			}
		case <-chexpire:
			// The timer has expired and chexpire has been closed.
		}

		// Callback is called unconditionally.  The caller should implement
		// cancel support.
		go callback()
		releaseChanId(chexpireId)
		timerDelete(t)
	}()

	return nil
}
