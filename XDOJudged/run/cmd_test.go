package run_test

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"syscall"
	"testing"
	"time"

	"linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/run"
)

func compareRE(expect *run.RuntimeError, get *run.RuntimeError) error {
	if expect == nil && get == nil {
		return nil
	}
	if expect != nil && get != nil {
		return nil
	}
	return fmt.Errorf("runtime error result mismatch: expect %v, get %v",
		expect, get)
}

func TestRuntimeError(t *testing.T) {
	type test struct{
		name string
		command string
		args []string
		attr *run.Attr
		sysattr *syscall.SysProcAttr
		re *run.RuntimeError
	}

	tl := time.Millisecond * 200
	idMap := []syscall.SysProcIDMap{
		{ContainerID: 0, HostID: os.Getuid(), Size: 1},
	}

	tests := []test{
		{name: "TestHelloWorld", command: "testdata/hw"},
		{
			name: "TestTLE",
			command: "testdata/loop",
			attr: &run.Attr{
				CPUTimeLimit: &tl,
			},
			re: &run.RuntimeError{
				Reason: run.ReasonCPUTimeLimit,
				Code: -int(unix.SIGKILL),
			},
		},
		{
			name: "TestILE",
			command: "testdata/loop",
			attr: &run.Attr{
				WallTimeLimit: &tl,
			},
			re: &run.RuntimeError{
				Reason: run.ReasonWallTimeLimit,
				Code: -int(unix.SIGKILL),
			},
		},
		{
			name: "TestNoCapability",
			command: "testdata/chroot",
			sysattr: &syscall.SysProcAttr{
				Chroot: "/",
				Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
				UidMappings: idMap,
				GidMappings: idMap,
			},
			re: &run.RuntimeError{
				Reason: run.ReasonUnknown,
				Code: 125,
			},
		},
		{name: "TestFork", command: "testdata/fork",},
		{
			name: "TestNoFork",
			command: "testdata/fork",
			attr: &run.Attr{
				CPUTimeLimit: &tl,
			},
			re: &run.RuntimeError{
				Reason: run.ReasonUnknown,
				Code: -int(unix.SIGSYS),
			},
		},
	}

	for _, i := range tests {
		t.Run(i.name, func(t *testing.T){
			cmd := run.Command(i.command, i.args...)
			cmd.Attr = i.attr
			cmd.SysProcAttr = i.sysattr
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			usage, re, err := cmd.Run()
			if err != nil {
				t.Fatal("can not run the command:", err)
			}
			t.Logf("usage = %v", usage)
			err = compareRE(i.re, re)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
