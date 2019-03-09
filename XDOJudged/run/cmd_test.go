package run_test

import (
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/bind"
	"linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/cgroup"
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
	case "TestHelloWorld", "TestChroot":
		fmt.Println("Hello, world.")
	case "TestILE":
		select {}
	case "TestTLE", "TestCgroup":
		for {
		}
	case "TestCapability", "TestNoCapability":
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
	t.Logf("usage = %v", usage)
	t.Logf("re = %v", re)
}

func TestRuntimeError(t *testing.T) {
	type transform func(*run.Cmd) interface{}
	type test struct {
		name      string
		attr      *run.Attr
		sysattr   *syscall.SysProcAttr
		re        *run.RuntimeError
		transform transform
	}

	tl := time.Millisecond * 200
	walltl := time.Millisecond * 400
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
				WallTimeLimit: &walltl,
			},
			re: &run.RuntimeError{
				Reason: run.ReasonWallTimeLimit,
				Code:   -int(unix.SIGKILL),
			},
		},
		{
			name: "TestCapability",
			attr: &run.Attr{KeepCap: true},
			sysattr: &syscall.SysProcAttr{
				Chroot:      "/",
				Cloneflags:  syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
				UidMappings: idMap,
				GidMappings: idMap,
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
			name:    "TestSigstop",
			sysattr: &syscall.SysProcAttr{Setpgid: true},
			re: &run.RuntimeError{
				Reason: run.ReasonUnknown,
				Code:   -int(unix.SIGSTOP),
			},
		},
		{
			// Still fragile.  Will fail on systems w/o /lib64, etc.
			name: "TestChroot",
			attr: &run.Attr{
				BindMount: []bind.BindMount{
					{
						OldDir: os.Args[0],
						NewDir: "/test",
					},
					{
						OldDir:   "/lib",
						NewDir:   "/lib",
						ReadOnly: true,
					},
					{
						OldDir:   "/lib64",
						NewDir:   "/lib64",
						ReadOnly: true,
					},
				},
			},
			sysattr: &syscall.SysProcAttr{
				Cloneflags:  unix.CLONE_NEWUSER | unix.CLONE_NEWNS,
				UidMappings: idMap,
				GidMappings: idMap,
			},
			transform: func(c *run.Cmd) interface{} {
				c.SysProcAttr.Chroot = filepath.Dir(c.Path)
				c.Path = "/test"
				return nil
			},
		},
		{
			name: "TestCgroup",
			attr: &run.Attr{
				CPUTimeLimit:  &tl,
				WallTimeLimit: &walltl,
			},
			transform: func(i *run.Cmd) interface{} {
				if os.Getenv("GO_XDOJ_RUN_TEST_CGROUP") != "1" {
					return "skip cgroup test by default"
				}
				cg, err := cgroup.Get(os.Getpid())
				if err != nil {
					return fmt.Errorf("can not get cgroup for test: %v",
						err)
				}
				err = cg.SetController([]cgroup.Controller{cgroup.CPU})
				if err != nil {
					return fmt.Errorf("SetController: %v", err)
				}
				_, err = cg.ToInnerNode()
				if err != nil {
					return fmt.Errorf("ToInnerNode: %v", err)
				}
				subcg, err := cg.Fork("test")
				if err != nil {
					return fmt.Errorf("cgroup Fork: %v", err)
				}
				// Throttle CPU to 10% so it would ILE
				i.Attr.Cgroup = subcg
				weightFile, err := subcg.OpenForWrite(cgroup.CPU, "max")
				if err != nil {
					return fmt.Errorf("OpenForWrite: %v", err)
				}
				defer weightFile.Close()
				_, err = weightFile.WriteString("10000 100000")
				if err != nil {
					return fmt.Errorf("Write: %v", err)
				}
				return nil
			},
			re: &run.RuntimeError{
				Reason: run.ReasonWallTimeLimit,
				Code:   -int(unix.SIGKILL),
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
			if i.transform != nil {
				result := i.transform(cmd)
				if result != nil {
					t.Skip(result)
				}
			}
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
