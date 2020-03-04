package proc

// BatchMemory reads from multiple addresses using a single syscall
// on backends which support it. On unsupported backends it will fall back to
// a syscall per address.
// When ReadMemory is called the address and buffer is stored but the memory
// requested is not actually read. Once all disperate reads are done you must call
// BatchMemory.BatchRead in order to actually load all the data and read the
// memory from the target process.
type BatchMemory map[uintptr][]byte

// ReadMemory will store the buffer passed in and associate it with the address given.
// Once BatchMemory.BatchRead is called the buffer will be filled with the data
// at the given address.
func (br BatchMemory) ReadMemory(buf []byte, addr uintptr) (int, error) {
	if _, ok := br[addr]; ok {
		return len(buf), nil
	}
	br[addr] = buf
	return len(buf), nil
}

// WriteMemory is here to satisfy the MemoryReadWriter interface, it does not actually do
// anything right now.
// TODO: could we benefit anywhere from batched writes?
func (br BatchMemory) WriteMemory(addr uintptr, data []byte) (written int, err error) {
	return len(data), nil
}

// BatchRead will attempt to read from all addresses requested in a single
// syscall. If the backend does not support this optimization then reads will
// be performed sequentially.
func (br BatchMemory) BatchRead(th Thread, mem MemoryReadWriter) error {
	return th.BatchRead(br)
}
