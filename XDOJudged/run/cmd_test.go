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
	if expect != nil && get != nil && *expect == *get {
		return nil
	}
	return fmt.Errorf("runtime error result mismatch: expect %v, get %v",
		expect, get)
}

func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_XDOJ_RUN_TEST_PROC") != "1" {
		return
	}
	defer os.Exit(0)

	testName := os.Args[2]
	switch testName {
	case "TestHelloWorld":
		fmt.Println("Hello, world.")
	case "TestILE":
		select {}
	case "TestTLE":
		for {
		}
	case "TestNoCapability":
		err := unix.Chroot("/")
		if err != nil {
			os.Exit(125)
		}
	case "TestFork", "TestNoFork":
		cmd := run.Command(os.Args[0], "-test.run=TestHelperProcess",
			"TestHelloWorld")
		cmd.Stdout = os.Stdout
		cmd.Start()
	}
}

func TestRuntimeError(t *testing.T) {
	type test struct {
		name    string
		command string
		attr    *run.Attr
		sysattr *syscall.SysProcAttr
		re      *run.RuntimeError
	}

	tl := time.Millisecond * 200
	idMap := []syscall.SysProcIDMap{
		{ContainerID: 0, HostID: os.Getuid(), Size: 1},
	}

	tests := []test{
		{name: "TestHelloWorld", command: "testdata/hw"},
		{
			name:    "TestTLE",
			command: "testdata/loop",
			attr: &run.Attr{
				CPUTimeLimit: &tl,
			},
			re: &run.RuntimeError{
				Reason: run.ReasonCPUTimeLimit,
				Code:   -int(unix.SIGKILL),
			},
		},
		{
			name:    "TestILE",
			command: "testdata/loop",
			attr: &run.Attr{
				WallTimeLimit: &tl,
			},
			re: &run.RuntimeError{
				Reason: run.ReasonWallTimeLimit,
				Code:   -int(unix.SIGKILL),
			},
		},
		{
			name:    "TestNoCapability",
			command: "testdata/chroot",
			sysattr: &syscall.SysProcAttr{
				Chroot:      "/",
				Cloneflags:  syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
				UidMappings: idMap,
				GidMappings: idMap,
			},
			re: &run.RuntimeError{
				Reason: run.ReasonUnknown,
				Code:   125,
			},
		},
		{name: "TestFork", command: "testdata/fork"},
		{
			name:    "TestNoFork",
			command: "testdata/fork",
			attr: &run.Attr{
				CPUTimeLimit: &tl,
			},
			re: &run.RuntimeError{
				Reason: run.ReasonUnknown,
				Code:   -int(unix.SIGSYS),
			},
		},
	}

	for _, i := range tests {
		t.Run(i.name, func(t *testing.T) {
			cmd := run.Command(os.Args[0], "-test.run=TestHelperProcess",
				i.name)
			cmd.Attr = i.attr
			cmd.SysProcAttr = i.sysattr
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			cmd.Env = []string{"GO_XDOJ_RUN_TEST_PROC=1"}
			usage, re, err := cmd.Run()
			if err != nil {
				t.Fatal("can not run the command:", err)
			}
			t.Logf("re = %v", re)
			t.Logf("usage = %v", usage)
			err = compareRE(i.re, re)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
