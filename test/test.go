package main

import (
	"fmt"
	"unsafe"

	"test/mpk"
)

func main() {
	pkey1, err := mpk.PkeyAlloc()
	if err != nil {
		fmt.Printf("Failed to allocate pkey, returned: %d\n", pkey1)
	} else {
		fmt.Printf("Allocated pkey: %d\n", pkey1)
	}

	pkey2, err := mpk.PkeyAlloc()
	if err != nil {
		fmt.Printf("Failed to allocate pkey, returned: %d\n", pkey2)
	} else {
		fmt.Printf("Allocated pkey: %d\n", pkey2)
	}

	var x int
	x = 1
	fmt.Println("Declaring var x with value:", x)

	pkru := mpk.AllRightsPKRU
	pkru = pkru.Update(pkey1, mpk.ProtX)
	pkru = pkru.Update(pkey2, mpk.ProtRX)
	mpk.WritePKRU(pkru)
	fmt.Println("Reading PKRU:", mpk.ReadPKRU())

	err = mpk.PkeyMprotect((uintptr(unsafe.Pointer(&x))>>12)<<12, 0x1000, pkey1)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Memory tagged with key 1")
	x = 2

	err = mpk.PkeyFree(pkey2)
	if err != nil {
		fmt.Println("Could not free pkey")
	} else {
		fmt.Println("pkey has been deallocated")
	}
}
