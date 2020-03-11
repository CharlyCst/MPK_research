package runtime

// tagPackages relies on sandboxes and pkgToId, they must be initialized before the call
func tagPackages() {
	println("[MPK] Beging tagging package")

	n := len(pkgToId)
	pkgAppearsIn := make(map[int][]int, n)

	for sbID, sb := range sandboxes {
		for _, pkg := range sb.Packages {
			// println(pkg) // Debug
			pkgID, ok := pkgToId[pkg]
			if !ok {
				println("Unable find package " + pkg)
				continue
			}

			sbGroup, ok := pkgAppearsIn[pkgID]
			if !ok {
				sbGroup = make([]int, 0)
			}
			pkgAppearsIn[pkgID] = append(sbGroup, sbID)
		}
	}

	sbKeys := make([][]int, len(sandboxes))
	for i := range sbKeys {
		sbKeys[i] = make([]int, 0)
	}

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
		// Add group key to sandbox
		for _, sbID := range SbGroupA {
			sbKeys[sbID] = append(sbKeys[sbID], key)
		}
		pkgGroups = append(pkgGroups, group)
	}

	println("[MPK] Done tagging packages")
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
