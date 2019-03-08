package bind

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
)

func (b *BindMount) DoMountWithChroot(chroot string) error {
	b, err := b.Sanitize()
	flag := unix.MS_BIND
	if !b.NoRecursive {
		flag |= unix.MS_REC
	}

	mountPoint := chroot + b.NewDir

	info, err := os.Lstat(b.OldDir)
	if err != nil {
		return fmt.Errorf("can not stat %s: %v", b.OldDir, err)
	}

	if info.IsDir() {
		err = os.MkdirAll(mountPoint, 0755)
	} else {
		info, err = os.Lstat(mountPoint)
		if os.IsNotExist(err) {
			var file *os.File
			file, err = os.Create(mountPoint)
			file.Close()
		}
	}

	if err != nil {
		return fmt.Errorf("can not create mount point %s: %v", b.NewDir,
			err)
	}

	err = unix.Mount(b.OldDir, mountPoint, "", uintptr(flag), "")
	if err != nil {
		return fmt.Errorf("can not bind mount %s to %s: %v", b.OldDir,
			mountPoint, err)
	}

	if b.ReadOnly {
		// modify the per-mount-point flags to be read-only
		err := unix.Mount(b.OldDir, mountPoint, "",
			unix.MS_BIND|unix.MS_REMOUNT|unix.MS_RDONLY, "")
		if err != nil {
			return fmt.Errorf("can not remount the bind mount %s "+
				"to be read only: %v", b.NewDir, err)
		}
	}
	return nil
}
