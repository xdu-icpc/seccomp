// Unit test for XDOJ Seccomp rule.
// Copyright (C) 2017-2018  Laboratory of ACM/ICPC, Xidian University

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

package seccomp_test

import (
	"bytes"
	"golang.org/x/net/bpf"
	"log"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"testing"

	"github.com/xdu-icpc/seccomp"
)

func init() {
	if len(os.Args) > 1 && os.Args[1] == "__child__" {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		ft, err := bpf.Assemble(noForkRule)
		if err != nil {
			log.Fatalf("bpf.Assemble: %v", err)
		}
		err = seccomp.SeccompFilter(0, ft)
		if err != nil {
			log.Fatalf("seccomp: %v", err)
		}
		err = syscall.Exec(os.Args[2], os.Args[3:], os.Environ())
		if err != nil {
			log.Fatalf("exec: %v", err)
		}
	}
}

func testRuleWith(t *testing.T, exe string, ret int) {
	self, err := os.Executable()
	if err != nil {
		t.Fatalf("can not get program name: %v", err)
	}

	var stderr bytes.Buffer
	cmd := exec.Command(self, "__child__", exe)
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err == nil {
		t.Fatal("the child returned 0")
	}

	exit, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatal(err)
	}

	t.Logf("child stderr: %v", stderr.String())
	sys, ok := exit.ProcessState.Sys().(syscall.WaitStatus)
	if !ok {
		t.Fatal("can not get wait status")
	}

	if !sys.Signaled() || sys.Signal() != syscall.SIGSYS {
		t.Fatalf("child status is %v instead of SIGSYS", sys)
	}
}

func TestRuleWithX86(t *testing.T) {
	if runtime.GOARCH != "386" && runtime.GOARCH != "amd64" {
		t.Skipf("GOARCH %s doesn't support this test", runtime.GOARCH)
	}

	testRuleWith(t, "testdata/x86", 2)
}

func TestRuleWithAMD64(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skipf("GOARCH %s doesn't support this test", runtime.GOARCH)
	}

	testRuleWith(t, "testdata/amd64", 2)
}

func TestRuleWithAMD64P32(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skipf("GOARCH %s doesn't support this test", runtime.GOARCH)
	}

	testRuleWith(t, "testdata/amd64p32", 2)
}
