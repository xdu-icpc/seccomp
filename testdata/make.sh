#!/bin/sh

if [ $(uname -m) = "x86_64" ]; then
	cc test.c -m64 -o amd64 -pthread
	cc test.c -m32 -o x86 -pthread
	cc test.c -mx32 -o amd64p32 -pthread
fi

if [ $(uname -m) = "i686" ]; then
	cc test.c -m32 -o x86 -pthread
fi
