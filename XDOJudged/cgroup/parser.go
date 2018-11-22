// Parser for cgroup information in /proc.
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

// +build linux

package cgroup

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Parse /proc/mounts to get cgroup mountpoints, return root Cgroup.
func parseProcMounts() (*Cgroup, error) {
	ret := Cgroup{}

	inf, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, err
	}
	defer inf.Close()

	rd := bufio.NewReader(inf)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		line = strings.TrimRight(line, "\n")
		fields := strings.Split(line, " ")

		mountpoint := fields[1]
		vfstype := fields[2]
		options := fields[3]
		var ctrls []Controller = nil

		switch vfstype {
		case "cgroup2":
			ctrls, err = getCgroup2Controllers(mountpoint)
		case "cgroup":
			ctrls = getControllersFromOption(options)
		}
		if err != nil {
			return nil, err
		}

		if ctrls != nil {
			ret.fs = append(ret.fs, mountpoint)
			for _, c := range ctrls {
				ret.fsid[c] = len(ret.fs)
			}
		}
	}
	return &ret, nil
}

func getCgroup2Controllers(path string) ([]Controller, error) {
	fn := fmt.Sprintf("%s/cgroup.controllers", path)
	inf, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer inf.Close()

	rd := bufio.NewReader(inf)
	line, err := rd.ReadString('\n')
	if err != nil && err != io.EOF {
		return nil, err
	}
	line = strings.TrimRight(line, "\n")
	names := strings.Split(line, " ")

	return append(getControllers(names), V2), nil
}

func getControllersFromOption(option string) []Controller {
	names := strings.Split(option, " ")
	return getControllers(names)
}

func getControllers(names []string) []Controller {
	ret := []Controller{}
	for _, s := range names {
		if c, ok := backMap[s]; ok {
			ret = append(ret, c)
		}
	}
	return ret
}

// Parse /proc/{pid}/cgroup to get the cgroups the process belongs to.
// path[i] would be the pathname relative to the mount point of
// the Controller i.  If Controller i is not enabled, path[i] would be
// empty.
//
// If multiple Controllers has the same mount point, only one of them
// will have non-empty path value.  For Controllers from cgroup v2, only
// V2 pseudo Controller will have non-empty path value.
//
// To distinguish from empty value, all pathes would begin with "/",
// just like the third field in /proc/{pid}/cgroup.
func parseProcCgroup(pid int) (path []string, err error) {
	fn := fmt.Sprintf("/proc/%d/cgroup", pid)
	inf, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer inf.Close()
	rd := bufio.NewReader(inf)

	path = make([]string, _CTRL_MAX)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		line = strings.TrimRight(line, "\n")
		tokens := strings.Split(line, ":")
		if tokens[1] == "" {
			path[V2] = tokens[2]
		} else {
			names := strings.Split(tokens[1], ",")
			if c, ok := backMap[names[0]]; ok {
				path[c] = tokens[2]
			}
		}
	}
	return
}
