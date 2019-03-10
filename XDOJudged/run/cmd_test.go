package run_test

import (
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/bind"
	"linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/run"
	"linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/run/testaux"
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

func compareRE(expect interface{}, get *run.RuntimeError) error {
	switch expect.(type) {
	case nil:
		if get == nil {
			return nil
		}
	case syscall.Signal:
		if get != nil && get.Signaled() && get.Signal() == expect {
			return nil
		}
	case int:
		if get.Exited() && get.ExitStatus() == expect {
			return nil
		}
	default:
		return fmt.Errorf("WTF is %v?", expect)
	}
	return fmt.Errorf("expected %v, get %v", expect, get)
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

func TestRun(t *testing.T) {
	type setup func(*run.Cmd) interface{}
	type test struct {
		name    string
		attr    *run.Attr
		sysattr *syscall.SysProcAttr
		expect  interface{}
		setup   setup
	}

	idMap := []syscall.SysProcIDMap{
		{ContainerID: 0, HostID: os.Getuid(), Size: 1},
	}

	tests := []test{
		{name: "TestHelloWorld"},
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
			expect: 125,
		},
		{name: "TestFork"},
		{
			name:   "TestNoFork",
			attr: &run.Attr{Seccomp: testaux.NoForkFilter},
			expect: syscall.SIGSYS,
		},
		{
			name:    "TestSigstop",
			sysattr: &syscall.SysProcAttr{Setpgid: true},
			expect:  syscall.SIGSTOP,
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
			setup: func(c *run.Cmd) interface{} {
				c.SysProcAttr.Chroot = filepath.Dir(c.Path)
				c.Path = "/test"
				return nil
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
			if i.setup != nil {
				result := i.setup(cmd)
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
			err = compareRE(i.expect, re)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
