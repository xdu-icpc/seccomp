//go:generate ./c2go.sh filter_c2go.go linux all
//go:generate gofmt -r "*_Ctype_struct_sock_filter -> *bpf.RawInstruction" -w zfilter_c2go_${GOOS}_${GOARCH}.go

package seccomp
