package run

import "golang.org/x/net/bpf"

func mustAssemble(insn []bpf.Instruction) (raw []bpf.RawInstruction) {
	raw, err := bpf.Assemble(insn)
	if err != nil {
		panic(err)
	}
	return
}

var noForkFilter = mustAssemble(noForkRule)
