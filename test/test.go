package main

import (
	"fmt"

	"test/mpk"
)

func main() {
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
	fmt.Printf("Reading PKRU: Ob%032b\n", mpk.ReadPKRU())

	err = mpk.PkeyFree(pkey)
	if err != nil {
		fmt.Println("Could not free pkey")
	} else {
		fmt.Println("pkey has been deallocated")
	}
}
