package testaux

import (
	"encoding/json"
	"golang.org/x/net/bpf"
)

func mustAssemble(insn []bpf.Instruction) (raw []bpf.RawInstruction) {
	raw, err := bpf.Assemble(insn)
	if err != nil {
		panic(err)
	}
	return
}

func mustEncodeToJson(insn []bpf.RawInstruction) string {
	b, err := json.Marshal(insn)
	if err != nil {
		panic(err)
	}
	return string(b)
}

var NoForkFilter = mustEncodeToJson(mustAssemble(noForkRule))
