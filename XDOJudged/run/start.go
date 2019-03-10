package run

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"io"
	"os"
)

var ErrBindWithoutChroot = errors.New("bind mount only makes sense with " +
	"chroot")
var ErrBindUnsafe = errors.New("bind mount is too dangerous with out new " +
	"mount namespace")

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

	// Maybe we want to keep the capabilites
	if attr.KeepCap {
		args = append(args, "-dropcap=false")
	}

	if attr.Seccomp != "" {
		args = append(args, "-seccomp=" + attr.Seccomp)
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
