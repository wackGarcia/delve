package proc

import (
	"syscall"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

type BatchMemoryReader map[uintptr][]byte

func (br BatchMemoryReader) ReadMemory(buf []byte, addr uintptr) (int, error) {
	if read, ok := br[addr]; ok {
		copy(buf, read)
		return len(buf), nil
	}
	br[addr] = buf
	return len(buf), nil
}

func (br BatchMemoryReader) WriteMemory(addr uintptr, data []byte) (written int, err error) {
	return len(data), nil
}

func (br BatchMemoryReader) BatchRead(tid int, mem MemoryReadWriter) error {
	_, err := ProcessVmReadBatch(tid, br)
	return err
	for addr, buf := range br {
		if _, err := mem.ReadMemory(buf, addr); err != nil {
			return err
		}
	}
	return nil //errors.New("not implemented")
}

func ProcessVmReadBatch(tid int, vecs map[uintptr][]byte) (int, error) {
	localvecs := make([]sys.Iovec, 0, len(vecs))
	remotevecs := make([]sys.Iovec, 0, len(vecs))
	for addr, buf := range vecs {
		len_iov := uint64(len(buf))
		local_iov := sys.Iovec{Base: &buf[0], Len: len_iov}
		remote_iov := sys.Iovec{Base: (*byte)(unsafe.Pointer(addr)), Len: len_iov}
		localvecs = append(localvecs, local_iov)
		remotevecs = append(remotevecs, remote_iov)
	}
	n, _, err := syscall.Syscall6(sys.SYS_PROCESS_VM_READV, uintptr(tid), uintptr(unsafe.Pointer(&localvecs[0])), uintptr(len(localvecs)), uintptr(unsafe.Pointer(&remotevecs[0])), uintptr(len(remotevecs)), 0)
	if err != syscall.Errno(0) {
		return 0, err
	}
	return int(n), nil
}
