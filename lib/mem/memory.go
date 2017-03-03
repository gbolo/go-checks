// mem provides an API for listing memory information on a unix type system
//
package mem

// Process is the generic interface that is implemented on every platform
// and provides common operations for processes.
type Memory interface {
	// Pid is the process ID for this process.
	Total() int64
	Free() int64
	Buffered() int64
	Cached() int64
	Available() int64
	UsedPercent() float64

	SwapTotal() int64
	SwapFree() int64
	SwapUsedPercent() float64
}

// Init return memory info.
//
// This of course will be a point-in-time snapshot of when this method was
// called. Some operating systems don't provide snapshot capability of the
// process table, in which case the process table returned might contain
// ephemeral entities that happened to be running when this was called.
func Init() (Memory, error) {
	return refreshMemory()
}
