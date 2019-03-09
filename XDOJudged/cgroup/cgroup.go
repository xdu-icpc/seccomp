// Cgroup basic functions
// Copyright (C) 2018-2019  Laboratory of ICPC, Xidian University

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

// +build linux

// package cgroup provides native Go interface to manipulate Linux
// cgroups(7) for resource control.
package cgroup

import (
	"errors"
	"fmt"
	"strings"
)

// Cgroup is a collection of processes that are bound to a set of defined
// limits or parameters.
type Cgroup struct {
	// Pathes to hierarchy
	fs []string

	// fs[fsid[i]-1] is controller i's hierarchy.
	// If fsid[i] = 0, controller i is not enabled.
	fsid [_CTRL_MAX]int

	// is this a inner node
	inner bool
}

func (cg *Cgroup) String() string {
	s := []string{}
	for c := Controller(0); c < _CTRL_MAX; c++ {
		if cg.fsid[c] != 0 {
			s = append(s, fmt.Sprintf("controller %v at %v", c,
				cg.fs[cg.fsid[c]-1]))
		}
	}
	return strings.Join(s, ", ")
}

// Get the Cgroup of a process with specified pid.
func Get(pid int) (*Cgroup, error) {
	cg, err := parseProcMounts()
	if err != nil {
		return nil, err
	}

	path, err := parseProcCgroup(pid)
	if err != nil {
		return nil, err
	}

	for c, p := range path {
		if p != "" && cg.fsid[c] != 0 {
			cg.fs[cg.fsid[c]-1] = cg.fs[cg.fsid[c]-1] + p
		}
	}

	// Some cgroup v2 controllers may be disabled in this subgroup.
	// Fix them up.
	if v2fsid := cg.fsid[V2]; v2fsid != 0 {
		active, err := getCgroup2Controllers(cg.fs[v2fsid-1])
		if err != nil {
			return nil, err
		}
		mark := [_CTRL_MAX]bool{}
		for _, c := range active {
			mark[c] = true
		}
		for c := Controller(0); c < _CTRL_MAX; c++ {
			if c != V2 && cg.fsid[c] == v2fsid && !mark[c] {
				cg.fsid[c] = 0
			}
		}
	}

	return cg, nil
}

var ErrNoController = errors.New("controller is not mounted")

// Set the Controllers used in the Cgroup.
//
// If one of the Controllers is not enabled in the Cgroup, it will return
// an error and do nothing.  Otherwise, it will configure the Cgroup to
// use ONLY the Controllers in the list.  It would be useful if some
// Controllers are enabled but we have no permission to use them.
func (cg *Cgroup) SetController(list []Controller) error {
	needv2 := false

	// check if the controllers in the list are enabled
	for _, c := range list {
		if v2, err := cg.IsV2(c); err == nil {
			needv2 = needv2 || v2
		} else {
			return err
		}
	}

	// throw unneeded hierarchies
	fs := make([]string, 0, len(cg.fs))
	done := make([]int, len(cg.fs))
	mark := make([]bool, _CTRL_MAX)
	for _, c := range list {
		fsidx := cg.fsid[c] - 1
		if done[fsidx] == 0 {
			fs = append(fs, cg.fs[fsidx])
			done[fsidx] = len(fs)
		}
		cg.fsid[c] = done[fsidx]
		mark[c] = true
	}
	if needv2 {
		mark[V2] = true
	}
	cg.fs = fs

	for c := Controller(0); c < _CTRL_MAX; c++ {
		if !mark[c] {
			cg.fsid[c] = 0
		}
	}

	return nil
}

// Check if a Controller is mounted as Cgroup V2.
func (cg *Cgroup) IsV2(c Controller) (bool, error) {
	if cg.fsid[c] == 0 {
		return false, ErrNoController
	}
	return cg.fsid[c] == cg.fsid[V2], nil
}

func writePid(path string, pid int) error {
	w, err := openForWrite(path + "/cgroup.procs")
	defer w.Close()

	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%d", pid)
	if err != nil {
		return err
	}
	return nil
}

// Attach the process with specified pid to the Cgroup.  If this function
// fails and returns an error, this process may be in inconsistent status
// and may have to be killed.
func (cg *Cgroup) Attach(pid int) error {
	// At first do permission check with cgroup v2 with correct delegation
	// support, if we have it.
	v2id := cg.fsid[V2] - 1
	if v2id != -1 {
		p := cg.fs[v2id]
		err := writePid(p, pid)
		if err != nil {
			// If permission check fails this would be EPERM.
			return err
		}
	}

	for i, p := range cg.fs {
		if i != v2id {
			err := writePid(p, pid)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
