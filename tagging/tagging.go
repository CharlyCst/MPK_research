package main

import (
	"fmt"
	"log"

	"github.com/KyleBanks/depth"
)

// Sandbox represents an isolated execution environment, with a name and a map of dependencies
type Sandbox struct {
	Name       string
	ID         int
	Deps       map[int]*Pkg // Induced dependencies
	directDeps map[int]*Pkg // Direct dependencies
	KeyGroups  []int
}

// NewSandbox returns a fresh sandbox
func NewSandbox(name string, id int, deps ...*Pkg) *Sandbox {
	directDeps := make(map[int]*Pkg, len(deps))
	for _, dep := range deps {
		directDeps[dep.ID] = dep
	}
	return &Sandbox{
		Name:       name,
		ID:         id,
		Deps:       make(map[int]*Pkg),
		directDeps: directDeps,
		KeyGroups:  make([]int, 0),
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
	return fmt.Sprintf("<Sandbox %s, ID %d, keys %v, deps [%s]>", sb.Name, sb.ID, sb.KeyGroups, deps)
}

// Pkg represents a package
type Pkg struct {
	Name string
	ID   int
}

// NewPkg initializes a fresh package
func NewPkg(name string, id int) *Pkg {
	return &Pkg{
		Name: name,
		ID:   id,
	}
}

func main() {
	// Initialize the global set of package, the strings package and a sandbox using strings
	pkgSet := make(map[int]*Pkg)     // !! Important: Packages must be zero indexed !!!
	sandboxes := make([]*Sandbox, 0) // !! Important: Sanboxes ID must be equal to their index in this array !!

	pkgIO := NewPkg("io", 0)
	pkgRuntime := NewPkg("runtime", 1)
	pkgSync := NewPkg("sync", 2)

	sandboxA := NewSandbox("sb_A", 0, pkgIO)
	sandboxB := NewSandbox("sb_B", 1, pkgRuntime)
	sandboxC := NewSandbox("sb_C", 2, pkgSync)

	sandboxes = append(sandboxes, sandboxA, sandboxB, sandboxC)

	pkgSet[0] = pkgIO
	pkgSet[1] = pkgRuntime
	pkgSet[2] = pkgSync

	crawlPackages("strings", pkgSet, sandboxes)
	tagPackages(pkgSet, sandboxes)

	fmt.Println()
	for _, pkg := range pkgSet {
		fmt.Printf("%3d %- 25s\n", pkg.ID, pkg.Name)
	}
	fmt.Println()
	for _, sb := range sandboxes {
		fmt.Println(sb)
	}
}

func tagPackages(pkgSet map[int]*Pkg, sandboxes []*Sandbox) {
	n := len(pkgSet)
	pkgAppearsIn := make(map[int][]int, n) // map packages to the list of sandboxes they appear in
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

	// KeyToSandboxes := make(map[int]map[int]struct{}) // Map key to set of sandboxes
	pkgGroups := make([][]int, 0)
	for len(pkgAppearsIn) > 0 {
		key := len(pkgGroups)
		group := make([]int, 0)
		_, SbGroupA := popMap(pkgAppearsIn)
		for id, SbGroupB := range pkgAppearsIn {
			if groupEqual(SbGroupA, SbGroupB) {
				group = append(group, id)
			}
		}
		for _, pkgID := range group {
			delete(pkgAppearsIn, pkgID)
		}
		for _, sbID := range SbGroupA {
			sandboxes[sbID].KeyGroups = append(sandboxes[sbID].KeyGroups, key)
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

func crawlPackages(rootPackage string, pkgSet map[int]*Pkg, sandboxes []*Sandbox) {
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
