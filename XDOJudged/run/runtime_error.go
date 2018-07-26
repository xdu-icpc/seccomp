package run

// RuntimeErrorReason represents a deduced reason of the runtime error.
type RuntimeErrorReason int

const (
	ReasonUnknown RuntimeErrorReason = iota
	ReasonWallTimeLimit
	ReasonCPUTimeLimit
	ReasonMemoryLimit
)

// A RuntimeError represents an error occured during external program
// (compiler, object code of submitted code, checker and etc.)  Note that
// the concept is wider than an "RE" in ACM/ICPC like contests.  Actually
// it also covers "MLE", "TLE", and "OLE".
type RuntimeError struct {
	// If positive, contains the exit code of the program.
	// If negative, contains (-s), where s is the signal killed the program.
	// If zero, there is a "race condition" - the program slightly violated
	// the limit, but managed to exit before to be killed.  A succeed run
	// of program should return nil, not a RuntimeError with Code 0.
	Code int

	// The deduced reason of runtime error.
	Reason RuntimeErrorReason
}
