package nid

import "slices"

// Sort the [NID]s.
func Sort(ids []NID) {
	slices.SortFunc(ids, Compare)
}

// SortBase the [Base] identifiers.
func SortBase(ids []Base) {
	slices.SortFunc(ids, CompareBase)
}
