// +build amd64 arm64 ppc64 s390x

package siginfo

type SiginfoHeader struct {
	Signo, Errno, Code int32
	pad0               int32
}
