package run

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"io"
	"os"
	"path/filepath"
	"time"

	"linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/posixtime"
)

var ErrBindWithoutChroot = errors.New("bind mount only makes sense with " +
	"chroot")
var ErrBindUnsafe = errors.New("bind mount is too dangerous with out new " +
	"mount namespace")

type errPathInvalid struct {
	path string
}

func (e errPathInvalid) Error() string {
	return fmt.Sprintf("path %s is invalid for bind", e.path)
}

func newErrPathInvalid(path string) errPathInvalid {
	return errPathInvalid{path: path}
}

func sanitizePathForBind(path string) (string, error) {
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	return "", newErrPathInvalid(path)
}

var zeroAttr Attr

func (c *Cmd) start() (err error) {
	attr := c.Attr
	if attr == nil {
		attr = &zeroAttr
	}

	self, err := os.Executable()
	if err != nil {
		return err
	}

	in, in1, err := os.Pipe()
	if err != nil {
		return err
	}
	defer in.Close()
	defer in1.Close()

	out1, out, err := os.Pipe()
	if err != nil {
		return err
	}
	defer out.Close()
	defer out1.Close()

	args := []string{ChildName}
	if attr.CPUTimeLimit == nil {
		// Do not use seccomp filter.
		args = append(args, "-seccomp=false")
	}

	// Delegate chroot to our helper.
	// If we chroot too early, we can't find this executable in chroot and
	// can not start the helper.
	chroot := ""
	if c.SysProcAttr != nil {
		chroot, c.SysProcAttr.Chroot = c.SysProcAttr.Chroot, chroot
	}
	if chroot != "" {
		args = append(args, "-chroot="+chroot)
	}

	if len(attr.BindMount) != 0 {
		if chroot == "" {
			return ErrBindWithoutChroot
		}
		if c.SysProcAttr.Cloneflags&unix.CLONE_NEWNS != unix.CLONE_NEWNS {
			return ErrBindUnsafe
		}
		for _, item := range attr.BindMount {
			item1, err := item.Sanitize()
			if err != nil {
				return err
			}
			args = append(args, "-bind="+item1.String())
		}
	}

	c.ExtraFiles = []*os.File{out1, in1}
	args = append(args, "--", c.Path)
	c.Args = append(args, c.Args...)
	c.Path = self

	err = c.Cmd.Start()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			c.kill()
			c.wait()
		}
	}()

	// Wait for the child, until it requests to Exec the real command
	in1.Close()
	var buf [256]byte
	_, err = in.Read(buf[:])
	if err != nil {
		return err
	}

	for i, b := range syncFlag {
		if buf[i] != b {
			return fmt.Errorf("child: %s", buf)
		}
	}

	// Setup the time limits
	clock, err := posixtime.GetCPUClockID(c.Process.Pid)
	if err != nil {
		return err
	}

	if attr.CPUTimeLimit != nil {
		c.cpuTimer, err = clock.AfterFunc(*c.Attr.CPUTimeLimit,
			func(ev posixtime.TimerEvent) {
				c.kill()
			})
		if err != nil {
			return err
		}
	}

	if attr.WallTimeLimit != nil {
		c.wallTimer = time.AfterFunc(*c.Attr.WallTimeLimit, func() {
			c.kill()
		})
	}

	// Set up resource limits
	for _, rlim := range attr.ResourceLimit {
		err := prlimit(c.Process.Pid, rlim.Resource, &rlim.Rlimit, nil)
		if err != nil {
			return err
		}
	}

	// Grant the child to continue
	out.Close()

	// Handle the error in child
	_, err = in.Read(buf[:])
	if err != io.EOF {
		return fmt.Errorf("child(pid = %d): %s", c.Process.Pid,
			bytes.NewBuffer(buf[:]).String())
	}

	return nil
}
