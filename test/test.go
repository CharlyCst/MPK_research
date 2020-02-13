package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"test/mpk"
)

func add(x, y int64) int64

func main() {
	println(add(3, 4))
	println(sysGetPID())
	pkey, err := mpk.PkeyAlloc()
	if err != nil {
		fmt.Printf("Failed to allocate pkey, returned: %d\n", pkey)
	} else {
		fmt.Printf("Allocated pkey: %d\n", pkey)
	}

	pkey, err = mpk.PkeyAlloc()
	if err != nil {
		fmt.Printf("Failed to allocate pkey, returned: %d\n", pkey)
	} else {
		fmt.Printf("Allocated pkey: %d\n", pkey)
	}

	mpk.WritePKRU(1<<2 + 1<<3) // set execute only on key 1
	fmt.Printf("Reading MKRU: %d\n", mpk.ReadPKRU())

	err = mpk.PkeyFree(pkey)
	if err != nil {
		fmt.Println("Could not free pkey")
	} else {
		fmt.Println("pkey has been deallocated")
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
