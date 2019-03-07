package run_test

import (
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"syscall"
	"testing"
	"time"

	"linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/run"
)

func init() {
	if os.Getenv("GO_XDOJ_RUN_TEST_PROC") == "1" {
		return
	}
	err := unix.Setrlimit(unix.RLIMIT_STACK, &unix.Rlimit{
		Cur: 8388608,
		Max: unix.RLIM_INFINITY,
	})
	if err != nil {
		panic(err)
	}
}

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
		if os.Getuid() != 0 {
			log.Fatalf("UID is not 0")
		}

		err := unix.Chroot("/")
		if err == unix.EPERM {
			os.Exit(125)
		}
	case "TestFork", "TestNoFork":
		cmd := run.Command(os.Args[0], "-test.run=TestHelperProcess",
			"TestHelloWorld")
		cmd.Stdout = os.Stdout
		cmd.Start()
	case "TestRlimit":
		rlim := unix.Rlimit{}
		err := unix.Getrlimit(unix.RLIMIT_STACK, &rlim)
		if err != nil {
			log.Fatalf("can not get rlimit: %v\n", err)
		}
		if rlim.Cur != unix.RLIM_INFINITY {
			log.Fatalf("stack is limited to %d\n", rlim.Cur)
		}
	case "TestSigstop":
		err := unix.Kill(0, unix.SIGSTOP)
		if err != nil {
			log.Fatalf("can not stop myself: %v\n", err)
		}
	}
}

func TestStart(t *testing.T) {
	cmd := run.Command(os.Args[0], "-test.run=TestHelperProcess", "TestTLE")
	cmd.Env = []string{"GO_XDOJ_RUN_TEST_PROC=1"}
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	cmd.Process.Kill()
	if err != nil {
		t.Fatal(err)
	}
	usage, re, err := cmd.Wait()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("usage = %v, re = %v", usage, re)
}

func TestRuntimeError(t *testing.T) {
	type test struct {
		name    string
		attr    *run.Attr
		sysattr *syscall.SysProcAttr
		re      *run.RuntimeError
	}

	tl := time.Millisecond * 200
	idMap := []syscall.SysProcIDMap{
		{ContainerID: 0, HostID: os.Getuid(), Size: 1},
	}

	tests := []test{
		{name: "TestHelloWorld"},
		{
			name: "TestTLE",
			attr: &run.Attr{
				CPUTimeLimit: &tl,
			},
			re: &run.RuntimeError{
				Reason: run.ReasonCPUTimeLimit,
				Code:   -int(unix.SIGKILL),
			},
		},
		{
			name: "TestILE",
			attr: &run.Attr{
				WallTimeLimit: &tl,
			},
			re: &run.RuntimeError{
				Reason: run.ReasonWallTimeLimit,
				Code:   -int(unix.SIGKILL),
			},
		},
		{
			name: "TestNoCapability",
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
		{name: "TestFork"},
		{
			name: "TestNoFork",
			attr: &run.Attr{
				CPUTimeLimit: &tl,
				ResourceLimit: []run.ResourceLimit{
					{
						Resource: unix.RLIMIT_CORE,
						Rlimit:   unix.Rlimit{Cur: 0, Max: 0},
					},
				},
			},
			re: &run.RuntimeError{
				Reason: run.ReasonUnknown,
				Code:   -int(unix.SIGSYS),
			},
		},
		{
			name: "TestRlimit",
			attr: &run.Attr{
				ResourceLimit: []run.ResourceLimit{
					{
						Resource: unix.RLIMIT_STACK,
						Rlimit: unix.Rlimit{
							Cur: unix.RLIM_INFINITY,
							Max: unix.RLIM_INFINITY,
						},
					},
				},
			},
		},
		{
			name: "TestSigstop",
			sysattr: &syscall.SysProcAttr{Setpgid: true},
			re: &run.RuntimeError{
				Reason: run.ReasonUnknown,
				Code:   -int(unix.SIGSTOP),
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
