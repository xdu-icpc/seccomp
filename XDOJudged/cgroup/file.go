// Utilities for cgroup fs files.
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
