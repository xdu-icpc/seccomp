// Cgroup controller definition and pretty printer.
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

package cgroup

// A Controller is a component that modifies the behavior of the processes
// in a cgroup.
type Controller int

const (
	V2 Controller = iota
	CPU
	CPUACCT
	CPUSET
	MEMORY
	DEVICES
	FREEZER
	NET_CLS
	BLKIO
	PERF_EVENT
	NET_PRIO
	HUGETLB
	PIDS
	IO
	_CTRL_MAX
)

var ctrlName = [_CTRL_MAX]string{
	"v2", "cpu", "cpuacct", "cpuset", "memory",
	"devices", "freezer", "net_cls", "blkio", "perf_event",
	"net_prio", "hugetlb", "pids", "io",
}

func getBackMap(name []string) map[string]Controller {
	ret := make(map[string]Controller)
	for i, s := range name {
		ret[s] = Controller(i)
	}
	return ret
}

var backMap = getBackMap(ctrlName[:])

func (c Controller) String() string {
	return ctrlName[c]
}
