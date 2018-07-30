package siginfo

type SiginfoHeader struct {
	Signo, Code, Errno int32
}
