package nid_test

import (
	"testing"

	"go.wamod.dev/nid"
)

func TestCompareBase(t *testing.T) {
	type args struct {
		a nid.Base
		b nid.Base
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "equal empty",
			args: args{
				a: nid.Base{},
				b: nid.Base{},
			},
			want: 0,
		},
		{
			name: "equal not empty",
			args: args{
				a: nid.Base{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				b: nid.Base{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			},
			want: 0,
		},
		{
			name: "less",
			args: args{
				a: nid.Base{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				b: nid.Base{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
			},
			want: -1,
		},
		{
			name: "greater",
			args: args{
				a: nid.Base{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
				b: nid.Base{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nid.CompareBase(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("CompareBase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	type args struct {
		a nid.NID
		b nid.NID
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "equal",
			args: args{
				a: nid.MustNaming("example").Apply(
					nid.Base{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				),
				b: nid.MustNaming("example").Apply(
					nid.Base{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				),
			},
			want: 0,
		},
		{
			name: "equal empty",
			args: args{
				a: nid.NID{},
				b: nid.NID{},
			},
			want: 0,
		},
		{
			name: "equal empty base",
			args: args{
				a: nid.MustNaming("example").Apply(nid.Base{1}),
				b: nid.MustNaming("example").Apply(nid.Base{1}),
			},
			want: 0,
		},
		{
			name: "less by name alpha",
			args: args{
				a: nid.MustNaming("example_a").Apply(nid.Base{1}),
				b: nid.MustNaming("example_b").Apply(nid.Base{1}),
			},
			want: -1,
		},
		{
			name: "greater by name alpha",
			args: args{
				a: nid.MustNaming("example_b").Apply(nid.Base{1}),
				b: nid.MustNaming("example_a").Apply(nid.Base{1}),
			},
			want: 1,
		},
		{
			name: "less by name numeric",
			args: args{
				a: nid.MustNaming("example_1").Apply(nid.Base{1}),
				b: nid.MustNaming("example_2").Apply(nid.Base{1}),
			},
			want: -1,
		},
		{
			name: "greater by name numeric",
			args: args{
				a: nid.MustNaming("example_2").Apply(nid.Base{1}),
				b: nid.MustNaming("example_1").Apply(nid.Base{1}),
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nid.Compare(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
