package nid_test

import (
	"testing"
	"time"

	"go.wamod.dev/nid"
)

func TestNewNaming(t *testing.T) {
	tt := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "short",
			input:   "example",
			wantErr: false,
		},
		{
			name:    "full",
			input:   "abcdefghijklmnopqrstuvwxyz_0123456789",
			wantErr: false,
		},
		{
			name:    "empty",
			input:   "",
			wantErr: true,
		},
		{
			name:    "starts underscore",
			input:   "_example",
			wantErr: true,
		},
		{
			name:    "ends underscore",
			input:   "example_",
			wantErr: true,
		},
		{
			name:    "double underscore",
			input:   "example__example",
			wantErr: true,
		},
		{
			name:    "starts number",
			input:   "123_example",
			wantErr: true,
		},
		{
			name:    "uppercase",
			input:   "EXAMPLE",
			wantErr: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := nid.NewNaming(tc.input)
			if tc.wantErr == (err == nil) {
				t.Errorf("NewNaming(); err = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

func TestMustNaming(t *testing.T) {
	tt := []struct {
		name      string
		input     string
		wantPanic bool
	}{
		{
			name:      "short",
			input:     "example",
			wantPanic: false,
		},
		{
			name:      "full",
			input:     "abcdefghijklmnopqrstuvwxyz_0123456789",
			wantPanic: false,
		},
		{
			name:      "empty",
			input:     "",
			wantPanic: true,
		},
		{
			name:      "starts underscore",
			input:     "_example",
			wantPanic: true,
		},
		{
			name:      "ends underscore",
			input:     "example_",
			wantPanic: true,
		},
		{
			name:      "double underscore",
			input:     "example__example",
			wantPanic: true,
		},
		{
			name:      "starts number",
			input:     "123_example",
			wantPanic: true,
		},
		{
			name:      "uppercase",
			input:     "EXAMPLE",
			wantPanic: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tc.wantPanic == (r == nil) {
					t.Errorf("MustNaming(); panic = %v, wantPanic = %v", r, tc.wantPanic)
				}
			}()
			nid.MustNaming(tc.input)
		})
	}
}

func TestNaming_New(t *testing.T) {
	tt := []struct {
		name      string
		idn       nid.Naming
		wantName  string
		wantPanic bool
	}{
		{
			name:      "not initialized",
			idn:       nid.Naming{},
			wantPanic: true,
		},
		{
			name:      "basic",
			idn:       nid.MustNaming("book"),
			wantName:  "book",
			wantPanic: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tc.wantPanic == (r == nil) {
					t.Errorf("Naming.New(); panic = %v; wantPanic = %v", r, tc.wantPanic)
				}
			}()

			if got := tc.idn.New(); got.Name() != tc.wantName {
				t.Errorf("Naming.New().Name() = %s; wantName = %s", got, tc.wantName)
			}
		})
	}
}

func TestNaming_NewAt(t *testing.T) {
	tt := []struct {
		name      string
		idn       nid.Naming
		time      time.Time
		wantName  string
		wantPanic bool
	}{
		{
			name:      "not initialized",
			idn:       nid.Naming{},
			time:      time.UnixMilli(12345),
			wantPanic: true,
		},
		{
			name:      "basic",
			idn:       nid.MustNaming("book"),
			time:      time.UnixMilli(12345),
			wantName:  "book",
			wantPanic: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tc.wantPanic == (r == nil) {
					t.Errorf("Naming.NewAt(); panic = %v; wantPanic = %v", r, tc.wantPanic)
				}
			}()

			got := tc.idn.NewAt(tc.time)

			if got.Name() != tc.wantName {
				t.Errorf("Naming.NewAt().Name() = %s; wantName = %s", got, tc.wantName)
			}

			if ts := got.Base().Time(); ts != tc.time {
				t.Errorf("Naming.NewAt().Base().Time() = %s; want = %s", ts, tc.time)
			}
		})
	}
}

func TestNaming_Is(t *testing.T) {
	tt := []struct {
		name      string
		idn       nid.Naming
		id        nid.NID
		want      bool
		wantPanic bool
	}{
		{
			name:      "not initialized",
			idn:       nid.Naming{},
			wantPanic: true,
		},
		{
			name: "same name",
			idn:  nid.MustNaming("book"),
			id:   nid.MustParse("book_000034o1ibe7u02570ak9evj9s"),
			want: true,
		},
		{
			name: "different name",
			idn:  nid.MustNaming("book"),
			id:   nid.MustParse("author_000034o1ibe7u02570ak9evj9s"),
			want: false,
		},
		{
			name: "empty",
			idn:  nid.MustNaming("book"),
			id:   nid.NID{},
			want: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tc.wantPanic == (r == nil) {
					t.Errorf("Naming.Is(); panic = %v; wantPanic = %v", r, tc.wantPanic)
				}
			}()

			if got := tc.idn.Is(tc.id); got != tc.want {
				t.Errorf("Naming.Is(id = %s) = %v, want %v", tc.id, got, tc.want)
			}
		})
	}
}

func TestNaming_Apply(t *testing.T) {
	tt := []struct {
		name      string
		idn       nid.Naming
		base      nid.Base
		want      nid.NID
		wantPanic bool
	}{
		{
			name:      "not initialized",
			idn:       nid.Naming{},
			wantPanic: true,
		},
		{
			name: "empty base",
			idn:  nid.MustNaming("book"),
			base: nid.Base{},
			want: nid.NID{},
		},
		{
			name: "non-empty base",
			idn:  nid.MustNaming("book"),
			base: nid.MustParseBase("000034o1ibe7u02570ak9evj9s"),
			want: nid.MustParse("book_000034o1ibe7u02570ak9evj9s"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tc.wantPanic == (r == nil) {
					t.Errorf("Naming.Apply(); panic = %v; wantPanic = %v", r, tc.wantPanic)
				}
			}()

			if got := tc.idn.Apply(tc.base); got != tc.want {
				t.Errorf("Naming.Apply(base = %s) = %v, want %v", tc.base, got, tc.want)
			}
		})
	}
}

func TestNaming_Update(t *testing.T) {
	tt := []struct {
		name      string
		idn       nid.Naming
		id        nid.NID
		want      nid.NID
		wantPanic bool
	}{
		{
			name:      "not initialized",
			idn:       nid.Naming{},
			wantPanic: true,
		},
		{
			name: "empty base",
			idn:  nid.MustNaming("book"),
			id:   nid.MustParse("author_00000000000000000000000000"),
			want: nid.MustParse("book_00000000000000000000000000"),
		},
		{
			name: "empty id",
			idn:  nid.MustNaming("book"),
			id:   nid.NID{},
			want: nid.MustParse("book_00000000000000000000000000"),
		},
		{
			name: "non-empty base",
			idn:  nid.MustNaming("book"),
			id:   nid.MustParse("author_000034o1ibe7u02570ak9evj9s"),
			want: nid.MustParse("book_000034o1ibe7u02570ak9evj9s"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tc.wantPanic == (r == nil) {
					t.Errorf("Naming.Update(); panic = %v; wantPanic = %v", r, tc.wantPanic)
				}
			}()

			if got := tc.idn.Update(tc.id); got != tc.want {
				t.Errorf("Naming.Update(id = %s) = %v, want %v", tc.id, got, tc.want)
			}
		})
	}
}
