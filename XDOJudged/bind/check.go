package bind

import (
	"fmt"
	"path/filepath"
)

type errPathInvalid struct{
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

func (b BindMount) Sanitize() (r *BindMount, err error) {
	b.OldDir, err = sanitizePathForBind(b.OldDir)
	if err != nil {
		return nil, err
	}
	b.NewDir, err = sanitizePathForBind(b.NewDir)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
