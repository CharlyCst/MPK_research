package mpk

import (
	"errors"
	"syscall"
)

// Syscall number on x86_64
const (
	sysPkeyMprotect = 329
	sysPkeyAlloc    = 330
	sysPkeyFree     = 331
)

func WritePKRU(prot uint32)
func ReadPKRU() uint32

// PkeyAlloc allocates a new pkey
func PkeyAlloc() (int, error) {
	pkey, _, _ := syscall.Syscall(sysPkeyAlloc, 0, 0, 0)
	if (int)(pkey) < 0 {
		return (int)(pkey), errors.New("Failled to allocate pkey")
	}
	return (int)(pkey), nil
}

// PkeyFree frees a pkey previously allocated
func PkeyFree(pkey int) error {
	result, _, _ := syscall.Syscall(sysPkeyFree, (uintptr)(pkey), 0, 0)
	if result != 0 {
		return errors.New("Could not free pkey")
	}
	return nil
}
