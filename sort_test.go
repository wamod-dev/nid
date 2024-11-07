package nid_test

import (
	"testing"

	"go.wamod.dev/nid"
)

func TestSort(t *testing.T) {
	list := []nid.NID{
		nid.MustNaming("example_b").Apply(nid.Base{1}),
		nid.MustNaming("example_2").Apply(nid.Base{1}),
		nid.MustNaming("example_a").Apply(nid.Base{1}),
		nid.MustNaming("example_1").Apply(nid.Base{1}),
		nid.MustNaming("example").Apply(nid.Base{1}),
		nid.MustNaming("example").Apply(nid.Base{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}),
		nid.MustNaming("example").Apply(nid.Base{1}),
		nid.MustNaming("example").Apply(nid.Base{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}),
		{},
	}

	want := []nid.NID{
		{},
		nid.MustNaming("example").Apply(nid.Base{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}),
		nid.MustNaming("example").Apply(nid.Base{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}),
		nid.MustNaming("example").Apply(nid.Base{1}),
		nid.MustNaming("example").Apply(nid.Base{1}),
		nid.MustNaming("example_1").Apply(nid.Base{1}),
		nid.MustNaming("example_2").Apply(nid.Base{1}),
		nid.MustNaming("example_a").Apply(nid.Base{1}),
		nid.MustNaming("example_b").Apply(nid.Base{1}),
	}

	nid.Sort(list)

	for i := 0; i < len(list); i++ {
		if list[i] != want[i] {
			t.Errorf("Sort()[%d] = %s, want = %s", i, list[i], want[i])
		}
	}
}

func TestSortBase(t *testing.T) {
	list := []nid.Base{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{},
	}

	want := []nid.Base{
		{},
		{},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
	}

	nid.SortBase(list)

	for i := 0; i < len(list); i++ {
		if list[i] != want[i] {
			t.Errorf("SortBase()[%d] = %s, want = %s", i, list[i], want[i])
		}
	}
}
