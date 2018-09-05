package cgroup

import (
	"os"
)

func (cg *Cgroup) getPath(c Controller, key string, prefixed bool) (string,
	error) {
	if cg.fsid[c] == 0 {
		return "", ErrNoController
	}
	if !prefixed {
		key = string(c) + "." + key
	}
	return cg.fs[cg.fsid[c]-1] + "/" + string(c) + "." + key, nil
}

func (cg *Cgroup) OpenForRead(c Controller, key string) (*os.File, error) {
	p, err := cg.getPath(c, key, false)
	if err != nil {
		return nil, err
	}
	return os.Open(p)
}

func openForWrite(p string) (*os.File, error) {
	return os.OpenFile(p, os.O_WRONLY, 0600)
}
