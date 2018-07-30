package siginfo

// Values for Siginfo.Code.  Positive values are reserved for
// kernel-generated signals.
const (
	// Sent by tkill.
	SI_TKILL = -6 + iota
	// Sent by queued SIGIO.
	SI_SIGIO
	// Sent by AIO completion.
	SI_ASYNCIO
	// Sent by real time mesq state change.
	SI_MESGQ
	// Sent by timer expiration.
	SI_TIMER
	// Sent by sigqueue.
	SI_QUEUE
	// Sent by kill, sigsend.
	SI_USER
	// Sent by asynch name lookup completion.
	SI_ASYNCNL = -60
	// Sent by kernel.
	SI_KERNEL = 0x80
)
