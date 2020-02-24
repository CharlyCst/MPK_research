package main

import (
	"fmt"
	"log"

	"github.com/KyleBanks/depth"
)

// Sandbox represents an isolated execution environment, with a name and a map of dependencies
type Sandbox struct {
	Name          string
	ID            int
	Deps          map[int]*Pkg // Complete dependencies
	directDeps    map[int]*Pkg // Direct dependencies
	depsToCluster []int
}

// NewSandbox returns a fresh sandbox
func NewSandbox(name string, id int, deps ...*Pkg) *Sandbox {
	directDeps := make(map[int]*Pkg, len(deps))
	for _, dep := range deps {
		directDeps[dep.ID] = dep
	}
	return &Sandbox{
		Name:          name,
		ID:            id,
		Deps:          make(map[int]*Pkg),
		directDeps:    directDeps,
		depsToCluster: make([]int, 0),
	}
}

func (sb *Sandbox) String() string {
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
	Name           string
	ID             int
	alwaysIncluded bool
	alwaysExcluded bool
	// usedIn         map[int]*Sandbox // Sandbox.ID set
}

// NewPkg initializes a fresh package
func NewPkg(name string, id int) *Pkg {
	return &Pkg{
		Name:           name,
		ID:             id,
		alwaysIncluded: true,
		alwaysExcluded: true,
		// usedIn:         make(map[int]*Sandbox),
	}
}

func main() {
	// Initialize the global set of package, the strings package and a sandbox using strings
	pkgSet := make(map[int]*Pkg) // !! Important: Packages must be zero indexed !!!

	pkgIO := NewPkg("io", 0)
	pkgRuntime := NewPkg("runtime", 1)
	pkgSync := NewPkg("sync", 2)

	sandboxA := NewSandbox("sb_A", 0, pkgIO)
	sandboxB := NewSandbox("sb_B", 1, pkgRuntime)
	sandboxC := NewSandbox("sb_C", 2, pkgSync)

	pkgSet[0] = pkgIO
	pkgSet[1] = pkgRuntime
	pkgSet[2] = pkgSync

	crawlPackages("strings", pkgSet, sandboxA, sandboxB, sandboxC)

	fmt.Println(sandboxA)
	fmt.Println(sandboxB)
	fmt.Println(sandboxC)

	tagPackages(pkgSet, sandboxA, sandboxB, sandboxC)

	fmt.Println()
	for _, pkg := range pkgSet {
		fmt.Printf("%2d %- 25s | included: %-5t | excluded: %-5t\n", pkg.ID, pkg.Name, pkg.alwaysIncluded, pkg.alwaysExcluded)
	}
}

func tagPackages(pkgSet map[int]*Pkg, sandboxes ...*Sandbox) {
	n := len(pkgSet)
	pkgAppearsIn := make(map[int][]int, n)
	for i := 0; i < n; i++ {
		pkgAppearsIn[i] = make([]int, 0)
	}
	for i := 0; i < n; i++ {
		for _, sb := range sandboxes {
			_, isInSb := sb.Deps[i]
			if isInSb {
				pkgAppearsIn[i] = append(pkgAppearsIn[i], sb.ID)
			}
		}
	}

	pkgGroups := make([][]int, 0)
	for len(pkgAppearsIn) > 0 {
		group := make([]int, 0)
		_, groupA := popMap(pkgAppearsIn)
		for id, groupB := range pkgAppearsIn {
			if groupEqual(groupA, groupB) {
				group = append(group, id)
			}
		}
		for _, id := range group {
			delete(pkgAppearsIn, id)
		}
		pkgGroups = append(pkgGroups, group)
	}
	fmt.Println(pkgGroups)
}

func groupEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func popMap(m map[int][]int) (int, []int) {
	for id, group := range m {
		return id, group
	}
	return -1, nil
}

func crawlPackages(rootPackage string, pkgSet map[int]*Pkg, sandboxes ...*Sandbox) {
	pkgID := len(pkgSet)
	pkgNameToID := make(map[string]int)
	pkgQueue := make([]depth.Pkg, 0)
	sbQueue := make([]int, 0)
	activeSb := make([]*Sandbox, 0)

	for id, pkg := range pkgSet {
		pkgNameToID[pkg.Name] = id
	}

	// Use depth to resolve dependency tree
	t := depth.Tree{
		ResolveInternal: true,
	}
	err := t.Resolve(rootPackage)
	if err != nil {
		log.Fatal(err)
	}

	pkgQueue = append(pkgQueue, *t.Root)
	sbQueue = append(sbQueue, 0)
	for len(pkgQueue) > 0 {
		lastIndex := len(pkgQueue) - 1
		pkg := pkgQueue[lastIndex]
		pkgQueue = pkgQueue[:lastIndex]
		nbSb := sbQueue[lastIndex]
		sbQueue = sbQueue[:lastIndex]

		// fmt.Println(pkg.Name)
		// for _, sb := range activeSb {
		// 	fmt.Print(sb.Name + " ")
		// }
		// fmt.Println()

		// Register the package if never seen before
		var pkgStruct *Pkg
		id, exist := pkgNameToID[pkg.Name]
		if !exist {
			pkgStruct = NewPkg(pkg.Name, pkgID)
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
			_, ok := sb.directDeps[id]
			if ok {
				nbNewSb++
				activeSb = append(activeSb, sb)
			}
		}

		// Register the current package into all active sandboxes
		for _, sb := range activeSb {
			sb.Deps[id] = pkgStruct
		}

		// Add all package dependencies to the queue
		for _, pkgDep := range pkg.Deps {
			pkgQueue = append(pkgQueue, pkgDep)
		}

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
