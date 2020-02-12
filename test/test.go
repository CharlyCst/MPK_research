package main

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

// Syscall number on x86_64
const (
	sysPkeyMprotect = 329
	sysPkeyAlloc    = 330
	sysPkeyFree     = 331
)

func add(x, y int64) int64

func main() {
	println(add(3, 4))
	println(sysGetPID())
	pkey, err := pkeyAlloc()
	if err != nil {
		fmt.Printf("Failed to allocate pkey, returned: %d\n", pkey)
	} else {
		fmt.Printf("Allocated pkey: %d\n", pkey)
	}
}

// Warning: doesn't work
func sysWrite(char uint) {
	syscall.Syscall(4, 1, 1, (uintptr)(unsafe.Pointer(&char)))
}

func sysGetPID() uint {
	pid, _, _ := syscall.Syscall(39, 0, 0, 0)
	return (uint)(pid)
}

func pkeyAlloc() (int, error) {
	pkey, _, _ := syscall.Syscall(sysPkeyAlloc, 0, 0, 0)
	if (int)(pkey) < 0 {
		return (int)(pkey), errors.New("Failled to allocate pkey")
	}
	return (int)(pkey), nil
}
