package cgroup

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
