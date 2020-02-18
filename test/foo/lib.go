package foo

import (
	"fmt"
)

type MpkgType struct {
	A [64]int
	B *int
}

var MpkgConst int

func Foo() {
		sandbox["",""] () {
			a := /*&MpkgType{}*/ new(int)
			fmt.Println("Sandboxed lib", a)
		}()
}