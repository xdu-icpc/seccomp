package cgroup_test

import (
	"os"
	"testing"

	"linux.xidian.edu.cn/git/XDU_ACM_ICPC/XDOJ-next/XDOJudged/cgroup"
)

func TestCgroup(t *testing.T) {
	cg, err := cgroup.Get(os.Getpid())
	if !t.Run("TestGet", func(t *testing.T) {
		if err != nil {
			t.Fatal(err)
		}
	}) {
		t.FailNow()
	}
	t.Logf("cgroup = %v", cg)

	err = cg.SetController([]cgroup.Controller{cgroup.MEMORY})
	if !t.Run("TestSetController", func(t *testing.T) {
		if err != nil {
			t.Fatal(err)
		}
	}) {
		t.FailNow()
	}
	t.Logf("cgroup = %v", cg)

	leaf, err := cg.ToInnerNode()
	if !t.Run("TestToInnerNode", func(t *testing.T) {
		if err != nil {
			t.Fatal(err)
		}
	}) {
		t.FailNow()
	}
	t.Logf("cgroup = %v", cg)
	t.Logf("leaf = %v", leaf)
}
