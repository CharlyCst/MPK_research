package main

import (
	"fmt"
	"unsafe"

	"test/mpk"
)

func main() {
	// gosb.Gosandbox()
	// x := sandbox["", ""] () {
	// 	foo.Foo()
	// }
	// // fmt.Println(x)

	// pkey, err := mpk.PkeyAlloc()
	// if err != nil {
	// 	fmt.Println("Could not allocate pkey")
	// 	return
	// }

	// pkru := mpk.AllRightsPKRU
	// mpk.WritePKRU(pkru)

	// fmt.Println()
	// x()
	// fmt.Println()
	// err = tagPackage("test/foo", pkey)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// pkru = pkru.Update(pkey, mpk.ProtX)
	// mpk.WritePKRU(pkru)

	// x()
	// fmt.Println()
	testMPK2()
}

// func tagPackage(packageName string, pkey mpk.Pkey) error {
// 	// id, ok := runtime.PkgToId()[packageName]
// 	// if !ok {
// 	// 	return errors.New("Could not find package")
// 	// }

// 	// for _, bloat := range runtime.PkgBloated() {
// 	// 	if bloat.Id == id {
// 	// 		// fmt.Println(bloat)
// 	// 		for _, pkgInfo := range bloat.Bloating.Relocs {
// 	// 			if pkgInfo.Addr != 0 && pkgInfo.Size != 0 {
// 	// 				fmt.Printf("%#x  %#x\n", pkgInfo.Addr, pkgInfo.Size)
// 	// 				err := mpk.PkeyMprotect(uintptr(pkgInfo.Addr), pkgInfo.Size, pkey)
// 	// 				if err != nil {
// 	// 					return errors.New("Could not mprotect package memory")
// 	// 				}
// 	// 			}
// 	// 		}
// 	// 	}
// 	// }

// 	// return nil
// }

type myStruct struct {
	myValue int
}

func testMPK2() {
	// Allocate an array
	a := make([]int, 1, 10000)

	// Allocate a key
	pkey, err := mpk.PkeyAlloc()
	check(err)

	// Tag the page containing `s` with our key
	err = mpk.PkeyMprotect(
		(uintptr(unsafe.Pointer(&a[0]))>>12)<<12, // Align pointer to page
		1<<12,                                // Page size
		mpk.SysProtRWX,                       // Base protection
		pkey,                                 // Key
	)

	fmt.Printf("Memory address of a[0]:  %p\n", a)
	fmt.Println("The value inside a[0]:  ", a[0])

	// Update a[0], then remove write access
	a[0] = 1

	pkru := mpk.AllRightsPKRU
	pkru = pkru.Update(pkey, mpk.ProtRX)
	mpk.WritePKRU(pkru)

	fmt.Println("The value inside a[0]:  ", a[0])
	a[0] = 2
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func testMPK() {
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

	err = mpk.PkeyMprotect((uintptr(unsafe.Pointer(&x))>>12)<<12, 0x1000, mpk.SysProtRWX, pkey1)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Memory tagged with key 1")
	x = 2

	pkru := mpk.AllRightsPKRU
	// pkru = pkru.Update(pkey1, mpk.ProtX)
	// pkru = pkru.Update(pkey2, mpk.ProtRX)
	mpk.WritePKRU(pkru)
	fmt.Println("Reading PKRU:", mpk.ReadPKRU())

	x = 3

	err = mpk.PkeyFree(pkey2)
	if err != nil {
		fmt.Println("Could not free pkey")
	} else {
		fmt.Println("pkey has been deallocated")
	}
}
