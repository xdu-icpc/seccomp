package run

import (
	"flag"
	"fmt"
	"golang.org/x/sys/unix"
	"io"
	"os"
	"runtime"

	"github.com/syndtr/gocapability/capability"
	"linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/seccomp"
)

// This command line (os.Args[0]) is internal used by the package.
// Don't clash.
const ChildName = "[xdoj child]"

var syncFlag = []byte{0x19, 0x26, 0x08, 0x17}

func closeOnExec(fd uintptr) error {
	_, _, errno := unix.RawSyscall(unix.SYS_FCNTL, fd, unix.F_SETFD,
		unix.FD_CLOEXEC)
	if errno != 0 {
		return errno
	}
	return nil
}

func bailOut(w io.Writer, msg string, err error) {
	pid := os.Getpid()
	fmt.Fprintf(w, "%s[%d]: %s: %v", ChildName, pid, msg, err)
	os.Exit(1)
}

func init() {
	if os.Args[0] == ChildName {
		cmdline := os.Args[1:]
		for i, arg := range os.Args {
			if arg == "--" {
				cmdline = os.Args[i+1:]
				break
			}
		}

		// We are in a child now.  Get the fds communicating with the
		// parent.
		in := os.NewFile(3, "|0")
		out := os.NewFile(4, "|1")

		err := closeOnExec(3)
		if err != nil {
			bailOut(out, "can not set |0 to close on exec", err)
		}
		err = closeOnExec(4)
		if err != nil {
			bailOut(out, "can not set |1 to close on exec", err)
		}

		// Parse argument list.
		fs := flag.NewFlagSet(ChildName, flag.ContinueOnError)

		useSeccomp := true
		fs.BoolVar(&useSeccomp, "seccomp", true, "enable seccomp filter")

		err = fs.Parse(os.Args[1:])
		if err != nil {
			bailOut(out, "can not parse arguments", err)
		}

		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		if useSeccomp {
			err = seccomp.SeccompFilter(0, noForkFilter)
			if err != nil {
				bailOut(out, "can not set seccomp filter", err)
			}
		}

		// TODO: Set up new namespace.

		capset, err := capability.NewPid2(0)

		if err != nil {
			bailOut(out, "can not init capability set", err)
		}

		capset.Clear(capability.BOUNDING)
		err = capset.Apply(capability.BOUNDING)
		if err != nil {
			bailOut(out, "can not clear capability bounding set", err)
		}

		out.Write(syncFlag[:])
		// Wait for the parent's permission for departure.
		b := make([]byte, 1)
		_, err = in.Read(b)
		if err != io.EOF {
			// Surrender the vessel immediately.
			os.Exit(1)
		}

		// We've set CLOEXEC.  So, if the parent get an EOF on the other
		// end of the pipe, it can be sure that execve() has succeeded.
		err = unix.Exec(cmdline[0], cmdline[1:], os.Environ())

		// Exec failed.  Dump the error message.
		bailOut(out, "exec failed", err)
	}
}
