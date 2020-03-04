package proc

import (
	"fmt"
)

type BatchMemoryReader map[uintptr][]byte

func (br BatchMemoryReader) ReadMemory(buf []byte, addr uintptr) (n int, err error) {
	fmt.Printf("BUF %#v\n", buf)
	br[addr] = buf
	return len(buf), nil
}

func (br BatchMemoryReader) WriteMemory(addr uintptr, data []byte) (written int, err error) {
	return len(data), nil
}

func (br BatchMemoryReader) BatchRead(mem MemoryReadWriter) error {
	for addr, buf := range br {
		if _, err := mem.ReadMemory(buf, addr); err != nil {
			return err
		}
	}
	fmt.Printf("%#v\n", br)
	return nil //errors.New("not implemented")
}
