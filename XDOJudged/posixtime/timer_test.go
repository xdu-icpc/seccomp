// Unitest of timer.go.
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

// +build linux

package posixtime_test

import (
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/posixtime"
)

func TestTimer(t *testing.T) {
	t0, err := posixtime.CLOCK_MONOTONIC.GetTime()
	if err != nil {
		t.Fatalf("can not get time of CLOCK_MONOTONIC: %v", err)
	}

	// To simplify the code we combine some tests in this function.
	// Use logs to distinguish them.
	t.Log("testing NewTimer...")

	timer, err := posixtime.CLOCK_MONOTONIC.NewTimer(sleepDuration)
	if err != nil {
		t.Fatalf("can not create timer: %v", err)
	}

	ev := <-timer.C
	t1 := ev.Time

	err = checkDuration(t1.Sub(*t0), sleepDuration)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("success.")

	t.Log("testing Reset...")
	err = timer.Reset(time.Millisecond * 200)
	if err != nil {
		t.Fatalf("can not reset timer: %v", err)
	}

	ev = <-timer.C
	t2 := ev.Time

	err = checkDuration(t2.Sub(t1), time.Millisecond*200)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("success.")

	t.Log("testing Stop active timer...")
	err = timer.Reset(time.Millisecond * 100)
	if err != nil {
		t.Fatalf("can not reset timer: %v", err)
	}

	if !timer.Stop() {
		t.Fatalf("returned false on active timer!")
	}
	gotimer := time.NewTimer(time.Millisecond * 200)
	select {
	case ev := <-timer.C:
		t.Fatalf("get unexpected value from C: %v", ev)
	case <-gotimer.C:
	}
	t.Log("success.")

	t.Log("testing Stop inactive timer...")
	timer.Reset(time.Millisecond * 10)
	time.Sleep(time.Millisecond * 100)
	if timer.Stop() {
		t.Fatalf("returned true on inactive timer!")
	}
	gotimer.Reset(time.Second)
	select {
	case <-gotimer.C:
		t.Fatalf("lost value in C!")
	case <-timer.C:
	}
	t.Log("success.")
}

func TestAfterFunc(t *testing.T) {
	t0 := time.Now()
	ch := make(chan error)
	posixtime.CLOCK_MONOTONIC.AfterFunc(time.Millisecond*200,
		func(ev posixtime.TimerEvent) {
			d := time.Since(t0)
			err := checkDuration(d, time.Millisecond*200)
			ch <- err
		})
	gotimer := time.NewTimer(time.Millisecond * 400)
	select {
	case <-gotimer.C:
		t.Fatalf("the func has not been executed at the expected time!")
	case err := <-ch:
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestIssue4(t *testing.T) {
	runtime.GOMAXPROCS(4)
	var timers [50]*posixtime.Timer
	// create many timers won't expire.
	for i := 0; i < 50; i++ {
		timer, err := posixtime.CLOCK_MONOTONIC.NewTimer(time.Hour)
		if err != nil {
			t.Fatalf("%v", err)
		}
		timers[i] = timer
	}

	// stop the timers in random sequence.
	src := rand.NewSource(19260817)
	rd := rand.New(src)
	perm := rd.Perm(50)
	for _, i := range perm {
		if ok := timers[i].Stop(); !ok {
			t.Fatalf("the timer is stopped unexpectedly.")
		}
	}

	// wait the timer goroutines to stop.
	time.Sleep(time.Millisecond * 100)

	// check if we leaked goroutines
	x := runtime.NumGoroutine()
	t.Logf("we have %d goroutines.", x)
	if x > 10 {
		t.Fatalf("goroutines leaked.")
	}
}

func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_POSIXTIME_HELPER_PROCESS") != "1" {
		return
	}

	// block forever
	ch := make(chan struct{})
	<-ch
}

func TestIssue14(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess")
	cmd.Env = []string{"GO_POSIXTIME_HELPER_PROCESS=1"}
	err := cmd.Start()
	if err != nil {
		t.Fatalf("can not start helper process: %v", err)
	}

	defer cmd.Wait()
	defer cmd.Process.Kill()

	clk, err := posixtime.GetCPUClockID(cmd.Process.Pid)
	if err != nil {
		t.Fatalf("GetCPUClockID: %v", err)
	}

	tm, err := clk.NewTimer(time.Second)
	if err != nil {
		t.Fatalf("NewTimer: %v", err)
	}

	// Kill and reap the process.
	err = cmd.Process.Kill()
	if err != nil {
		t.Fatalf("Kill: %v", err)
	}

	err = cmd.Wait()
	_, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("Wait: %v", err)
	}

	// Try to stop the timer on the reaped process
	armed := tm.Stop()
	if !armed {
		t.Fatalf("Stop: wrong return value")
	}
}
