package nid

// Compare two [NID] identifiers.
func Compare(a, b NID) int {
	if a.name < b.name {
		return -1
	} else if a.name > b.name {
		return 1
	}

	return CompareBase(a.base, b.base)
}

// CompareBase two [Base] identifiers
func CompareBase(a, b Base) int {
	for i := 0; i < baseLen; i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}

	return 0
}
