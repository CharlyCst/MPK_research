package main

import (
	"fmt"
	"log"

	"github.com/KyleBanks/depth"
)

// Sandbox represents an isolated execution environment, with a name and a map of dependencies
type Sandbox struct {
	Name string
	ID   int
	Deps map[int]*Pkg // Pkg.ID set
}

func (sb Sandbox) String() string {
	deps := ""
	i := 0
	for _, dep := range sb.Deps {
		if i > 0 {
			deps += ", "
		}
		deps += dep.Name
		i++
	}
	return fmt.Sprintf("<Sandbox %s, ID %d, deps [%s]>", sb.Name, sb.ID, deps)
}

// Pkg represents a package
type Pkg struct {
	Name   string
	ID     int
	usedIn map[int]*Sandbox // Sandbox.ID set
}

func main() {
	// Initialize the global set of package, the strings package and a sandbox using strings
	pkgSet := make(map[int]*Pkg)
	pkgIO := Pkg{
		Name:   "io",
		ID:     0,
		usedIn: make(map[int]*Sandbox),
	}
	sandboxA := Sandbox{
		Name: "sb_A",
		ID:   0,
		Deps: make(map[int]*Pkg),
	}

	pkgSet[0] = &pkgIO
	sandboxA.Deps[0] = &pkgIO
	pkgIO.usedIn[0] = &sandboxA

	crawlPackage("strings", pkgSet, &sandboxA)
	// fmt.Println(pkgSet)
	fmt.Println(sandboxA)
}

func crawlPackage(rootPackage string, pkgSet map[int]*Pkg, sandboxes ...*Sandbox) {
	pkgID := 0
	pkgNameToID := make(map[string]int)
	pkgQueue := make([]depth.Pkg, 0)
	sbQueue := make([]int, 0)
	activeSb := make([]*Sandbox, 0)

	for id, pkg := range pkgSet {
		if id <= pkgID {
			pkgID = id + 1
		}

		pkgNameToID[pkg.Name] = id
	}

	t := depth.Tree{
		ResolveInternal: true,
	}
	err := t.Resolve("strings")
	if err != nil {
		log.Fatal(err)
	}

	for _, dep := range t.Root.Deps {
		pkgQueue := append(pkgQueue, dep)
		sbQueue := append(sbQueue, 0)

		for len(pkgQueue) > 0 {
			lastIndex := len(pkgQueue) - 1
			pkg := pkgQueue[lastIndex]
			pkgQueue = pkgQueue[:lastIndex]
			nbSb := sbQueue[lastIndex]
			sbQueue = sbQueue[:lastIndex]

			// Register the package if never seen before
			var pkgStruct *Pkg
			id, exist := pkgNameToID[pkg.Name]
			if !exist {
				pkgStruct = &Pkg{
					Name:   pkg.Name,
					ID:     pkgID,
					usedIn: make(map[int]*Sandbox),
				}
				pkgSet[pkgID] = pkgStruct
				pkgNameToID[pkg.Name] = pkgID
				id = pkgID
				pkgID++
			} else {
				pkgStruct, _ = pkgSet[id]
			}

			// Add newly activated sandboxes
			nbNewSb := 0
			for _, sb := range sandboxes {
				_, ok := sb.Deps[id]
				if ok {
					nbNewSb++
					activeSb = append(activeSb, sb)
				}
			}

			// Register the current package into all active sandboxes
			for _, sb := range activeSb {
				sb.Deps[id] = pkgStruct
			}

			// fmt.Print(pkg.Name, " | ")
			for _, pkgDep := range pkg.Deps {
				// fmt.Print(" ", pkgDep.Name)
				pkgQueue = append(pkgQueue, pkgDep)
			}
			// fmt.Println()

			// Remove out of scope sandboxes
			if len(pkg.Deps) == 0 {
				activeSb = activeSb[:len(activeSb)-nbSb-nbNewSb]
			} else {
				sbQueue = append(sbQueue, nbSb+nbNewSb)
				for i := 0; i < len(pkg.Deps)-1; i++ {
					sbQueue = append(sbQueue, 0)
				}
			}
		}
	}
}
