// High level package for XDOJ sandbox function.
package run

import (
	"flag"
	"fmt"
	"io"
	"golang.org/x/sys/unix"
	"os"
	"runtime"

	"linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/seccomp"
	"github.com/syndtr/gocapability/capability"
)

// This command line (os.Args[0]) is internal used by the package.
// Don't clash.
const ChildName = "[xdoj child]"

var syncFlag = []byte{0x19, 0x26, 0x08, 0x17};

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

		// Parse argument list.
		fs := flag.NewFlagSet(ChildName, flag.ContinueOnError)

		useSeccomp := true
		fs.BoolVar(&useSeccomp, "seccomp", true, "enable seccomp filter")

		err := fs.Parse(os.Args[1:])
		if err != nil {
			bailOut(out, "can not parse arguments", err)
		}

		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		if useSeccomp {
			err = seccomp.SeccompFilter(0, seccomp.NoForkFilter)
			if err != nil {
				bailOut(out, "can not set seccomp filter", err)
			}
		}

		// TODO: Set up new namespace.
		cap, err := capability.NewPid2(0)

		if err != nil {
			bailOut(out, "can not init capability set", err)
		}

		cap.Clear(capability.BOUNDING)
		err = cap.Apply(capability.BOUNDING)
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

		// The Go pipes have CLOEXEC set.  We let execve() to close |1,
		// instead of close it manually.  So, if the parent get an EOF
		// on the other end of the pipe, it can be sure that execve() has
		// succeeded.
		err = unix.Exec(cmdline[0], cmdline[1:], os.Environ())

		// Exec failed.  Dump the error message.
		bailOut(out, "exec failed", err)
	}
}
