// Demultiplex POSIX timer's signals.
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
	"sync/atomic"
	"syscall"

	"linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/siginfo"
)

func demux() {
	var info siginfo.Siginfo

	mask := sigsetRTMIN()
	for {
		err := sigwaitinfo(&mask, &info)
		if err == syscall.EINTR {
			// If unlucky a signal other than SIGRTMIN may be delivered
			// to our thread.  In this case just do next loop.
			continue
		}
		if err != nil {
			panic(err)
		}

		id, queued := -1, false
		if p := info.ToSiginfoTimer(); p != nil {
			id = int(p.Value)
		} else if p := info.ToSiginfoQueue(); p != nil {
			id = int(p.Value)
			queued = true
		}
		if id == -1 {
			continue
		}

		ch := getChanById(id)
		esrch := atomic.SwapInt32(&esrchFlag[id], 1)
		if esrch == 1 && !queued {
			// The process has been reaped and timerRoutine deliberately
			// queued a signal.  Ignore other signals before it.
			continue
		}

		close(ch)
	}
}
