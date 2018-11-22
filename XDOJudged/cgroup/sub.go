// Sub-cgroup creation.
// Copyright (C) 2018  Laboratory of ICPC, Xidian University

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Author: Xi Ruoyao <xry111@mengyan1223.wang>

package cgroup

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

func (cg *Cgroup) fork(name string) (*Cgroup, error) {
	fs := make([]string, 0, len(cg.fs))
	for _, path := range cg.fs {
		p := path + "/" + name
		err := os.Mkdir(p, 0700)
		if err != nil {
			// roll back
			for _, p := range fs {
				os.Remove(p)
			}
			return nil, err
		}
		fs = append(fs, p)
	}

	return &Cgroup{fs: fs, fsid: cg.fsid}, nil
}

func (cg *Cgroup) pushdownV2() error {
	if cg.fsid[V2] == 0 {
		// nothing to do
		return nil
	}

	f, err := openForWrite(cg.fs[cg.fsid[V2]-1] +
		"/cgroup.subtree_control")
	if err != nil {
		return err
	}
	defer f.Close()

	for c := Controller(0); c < _CTRL_MAX; c++ {
		if c != V2 && cg.fsid[c] == cg.fsid[V2] {
			_, err := fmt.Fprintf(f, "+%s", c)
			if err != nil {
				return err
			}
		}
	}
	f.Close()
	return nil
}

// ToInnerNode converts a Cgroup to an inner node in the cgroup hierarchy.
// It creates an leaf node under the Cgroup cg, and move all processes
// under it into the leaf.
//
// If ToInnerNode fails and return an error, it's likely to leave an
// inconsistent cgroup hierarchy.  Maybe panic is the only thing we can do.
func (cg *Cgroup) ToInnerNode() (leaf *Cgroup, err error) {
	leaf, err = cg.fork("_leaf_")
	if err != nil {
		return
	}

	// Now move all processes into the leaf node.
	for i, path := range cg.fs {
		set := make(map[int]struct{})
		fn := path + "/cgroup.procs"
		for done := false; !done; {
			done = true

			f, err := os.Open(fn)
			if err != nil {
				return nil, err
			}

			wf, err := openForWrite(leaf.fs[i] + "/cgroup.procs")
			if err != nil {
				f.Close()
				return nil, err
			}

			pid := 0
			for {
				n, err := fmt.Fscanf(f, "%d", &pid)
				if err == io.EOF {
					break
				}

				if err != nil {
					wf.Close()
					f.Close()
					return nil, err
				}

				if n != 1 {
					wf.Close()
					f.Close()
					return nil, io.ErrNoProgress
				}

				if _, ok := set[pid]; ok {
					// cgroup.procs may contain duplicated pids
					continue
				}

				if err != nil && err != syscall.ESRCH {
					wf.Close()
					f.Close()
					return nil, err
				}

				fmt.Fprintf(wf, "%d", pid)

				done = false
				set[pid] = struct{}{}
			}
		}
	}
	err = cg.pushdownV2()
	if err != nil {
		return nil, err
	}

	return leaf, nil
}

func (cg *Cgroup) Fork(name string) (child *Cgroup, err error) {
	child, err = cg.fork(name)
	if err != nil {
		return nil, err
	}

	err = cg.pushdownV2()
	if err != nil {
		return nil, err
	}

	return
}
