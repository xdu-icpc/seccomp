#!/bin/bash

# A wrapper for `go tool cgo -godefs`, to generate architecture specific
# file with information in system headers.

# Copyright (C) 2017  Laboratory of ACM/ICPC, Xidian University

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

# Author: Xi Ruoyao <ryxi@stu.xidian.edu.cn>

# This just works.  DO NOT use this in your project.

if [ $# -lt 3 ]; then
	echo "usage: $0 {input} {oslist} {archlist}"
	exit 2
fi

inputfile="$1"
outputfile=$(echo $inputfile | sed "s/_${GOOS}\\.go/.go/; s/\\.go/_${GOOS}_${GOARCH}.go/; s/^.*$/z&/")
oslist="$2"
archlist="$3"

osok=0
if [ "$oslist" = "all" ]; then
	osok=1
elif echo "$oslist" | grep -q "$GOOS"; then
	osok=1
fi

platok=0
if [ "$archlist" = "all" ]; then
	platok=1
elif echo "$archlist" | grep -q "$GOARCH"; then
	platok=1
fi

if [ $osok -a $platok ]; then
	echo "$inputfile -> $outputfile"
	go tool cgo -godefs $inputfile > $outputfile
	if [ ! $? ]; then
		rm $outputfile
	fi
	rm -rf _obj
	gofmt -w $outputfile
fi
