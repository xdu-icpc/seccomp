// Cgroup unit tests.  Shoule be run with run_test.sh.
// Copyright (C) 2018  Laboratory of ICPC, Xidian University

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

// Author: Xi Ruoyao <xry111@mengyan1223.wang>

package cgroup_test

import (
	"os"
	"os/exec"
	"os/signal"
	"testing"

	"linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/cgroup"
)

func TestHelperProcess(t *testing.T) {
	if os.Getenv("XDOJ_CGROUP_WANT_TEST_HELPER") != "1" {
		t.Skip()
	}
	defer os.Exit(0)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
}

func TestCgroup(t *testing.T) {
	cg, err := cgroup.Get(os.Getpid())
	if !t.Run("TestGet", func(t *testing.T) {
		if err != nil {
			t.Fatal(err)
		}
	}) {
		t.FailNow()
	}
	t.Logf("cgroup = %v", cg)

	err = cg.SetController([]cgroup.Controller{cgroup.MEMORY})
	if !t.Run("TestSetController", func(t *testing.T) {
		if err != nil {
			t.Fatal(err)
		}
	}) {
		t.FailNow()
	}
	t.Logf("cgroup = %v", cg)

	leaf, err := cg.ToInnerNode()
	if !t.Run("TestToInnerNode", func(t *testing.T) {
		if err != nil {
			t.Fatal(err)
		}
	}) {
		t.FailNow()
	}
	t.Logf("cgroup = %v", cg)
	t.Logf("leaf = %v", leaf)

	testcg, err := cg.Fork("test")
	if !t.Run("TestFork", func(t *testing.T) {
		if err != nil {
			t.Fatal(err)
		}
	}) {
		t.FailNow()
	}
	t.Logf("testcg = %v", testcg)

	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess")
	cmd.Env = []string{"XDOJ_CGROUP_WANT_TEST_HELPER=1"}
	err = cmd.Start()
	if err != nil {
		t.Fatal(err)
	}

	childcg, err := cgroup.Get(cmd.Process.Pid)
	if !t.Run("TestGetChildCgroup", func(t *testing.T) {
		if err != nil {
			t.Fatal(err)
		}
	}) {
		t.FailNow()
	}
	t.Logf("child cgroup = %v", childcg)

	err = testcg.Attach(cmd.Process.Pid)
	if !t.Run("TestAttach", func(t *testing.T) {
		if err != nil {
			t.Fatal(err)
		}
		childcg, err := cgroup.Get(cmd.Process.Pid)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("child cgroup changed to: %v", childcg)
	}) {
		t.FailNow()
	}

	cmd.Process.Signal(os.Interrupt)
}
