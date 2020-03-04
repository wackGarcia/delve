// +build linux,amd64 linux,arm64

package native

import (
	"syscall"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

// Iovec is copied from golang.org/x/sys/unix with one modification.
// In the original struct, Base was of type *byte. This causes an issue
// for us as we store a pointer outside of our address space which is also
// not known to the Go GC. This would cause the garbage collector to panic
// if a sweep happened while we held a *byte pointer to an unknown address.
type Iovec struct {
	Base uintptr
	Len  uint64
}

// ProcessVmRead calls process_vm_readv.
func ProcessVmRead(tid int, addr uintptr, data []byte) (int, error) {
	len_iov := uint64(len(data))
	local_iov := Iovec{Base: uintptr(unsafe.Pointer(&data[0])), Len: len_iov}
	remote_iov := Iovec{Base: uintptr(unsafe.Pointer(addr)), Len: len_iov}
	p_local := uintptr(unsafe.Pointer(&local_iov))
	p_remote := uintptr(unsafe.Pointer(&remote_iov))
	n, _, err := syscall.Syscall6(sys.SYS_PROCESS_VM_READV, uintptr(tid), p_local, 1, p_remote, 1, 0)
	if err != syscall.Errno(0) {
		return 0, err
	}
	return int(n), nil
}

// ProcessVmWrite calls process_vm_writev.
func ProcessVmWrite(tid int, addr uintptr, data []byte) (int, error) {
	len_iov := uint64(len(data))
	local_iov := Iovec{Base: uintptr(unsafe.Pointer(&data[0])), Len: len_iov}
	remote_iov := Iovec{Base: uintptr(unsafe.Pointer(addr)), Len: len_iov}
	p_local := uintptr(unsafe.Pointer(&local_iov))
	p_remote := uintptr(unsafe.Pointer(&remote_iov))
	n, _, err := syscall.Syscall6(sys.SYS_PROCESS_VM_WRITEV, uintptr(tid), p_local, 1, p_remote, 1, 0)
	if err != syscall.Errno(0) {
		return 0, err
	}
	return int(n), nil
}

// ProcessVmReadBatch reads data from multiple addresses in a single syscall.
func ProcessVmReadBatch(tid int, vecs map[uintptr][]byte) (int, error) {
	localvecs := make([]Iovec, 0, 10)
	remotevecs := make([]Iovec, 0, 10)
	for addr, buf := range vecs {
		len_iov := uint64(len(buf))
		localvecs = append(localvecs, Iovec{Base: uintptr(unsafe.Pointer(&buf[0])), Len: len_iov})
		remotevecs = append(remotevecs, Iovec{Base: uintptr(unsafe.Pointer(addr)), Len: len_iov})
	}
	n, _, err := syscall.Syscall6(
		sys.SYS_PROCESS_VM_READV,
		uintptr(tid), uintptr(unsafe.Pointer(&localvecs[0])),
		uintptr(len(localvecs)), uintptr(unsafe.Pointer(&remotevecs[0])),
		uintptr(len(remotevecs)), 0)
	if err != syscall.Errno(0) {
		return 0, err
	}
	return int(n), nil
}
