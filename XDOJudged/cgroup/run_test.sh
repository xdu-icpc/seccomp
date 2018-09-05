#!/bin/sh

# Run test of this package with systemd-run
systemd-run -t -p Delegate=yes -p DynamicUser=yes -E GOPATH=$GOPATH -- `which go` test linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/cgroup -v
