package siginfo

type SiginfoHeader struct {
	Signo, Code, Errno int32
	pad0               int32
}
