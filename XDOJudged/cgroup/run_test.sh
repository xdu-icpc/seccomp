#!/bin/sh

# Run test of this package with systemd-run.
# Copyright (C) 2018  Laboratory of ICPC, Xidian University

# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published
# by the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.

# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

# Author: Xi Ruoyao <xry111@mengyan1223.wang>

systemd-run -t -p Delegate=yes -p DynamicUser=yes -E GOPATH=$GOPATH -- `which go` test linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/cgroup -v
