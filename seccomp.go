// Document of package seccomp.
// Copyright (C) 2017  Laboratory of ACM/ICPC, Xidian University

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

// Author: Xi Ruoyao <ryxi@stu.xidian.edu.cn>

// Package seccomp contains Seccomp BPF filters and support routines.
//
// Seccomp is highly dependant on architecture.  It had cost months
// to implement a generic Seccomp package but failed.  So this package
// JUST work - do NOT use it outside XDOJ unless you know what you are
// doing.
//
// This package assume you'd like to use Seccomp filters.  So it use
// prctl(2) to set no_new_privs bit at startup (even before Go runtime
// initialization).  If you have to run setuid programs, unset it, but
// reset it before installing filter.
//
// XDOJ should only use NoForkFilter and use other techniques to limit
// system resources other than PGID.  But before Linux 4.8 there was an
// issue so we have to filter out ptrace syscall.
package seccomp
