package nid_test

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"reflect"
	"testing"

	"go.wamod.dev/nid"
)

// Example demonstrates how to use the [go.wamod.io/nid] package.
func Example() {
	// Create a new [nid.Naming] for books
	bookIDN := nid.MustNaming("book")

	// Create one new [nid.NID] for
	bookID := bookIDN.New()

	fmt.Printf("Book ID: %s", bookID)
}

func TestMustParse(t *testing.T) {
	tt := []struct { //nolint:dupl
		name      string
		str       string
		want      nid.NID
		wantPanic bool
	}{
		{
			name: "valid",
			str:  "book_000034o1ibe7u02570ak9evj9s",
			want: nid.MustNaming("book").Apply(
				nid.Base{
					0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
					0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
				},
			),
		},
		{
			name: "zeros",
			str:  "book_00000000000000000000000000",
			want: nid.MustNaming("book").Apply(
				nid.Base{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
			),
		},
		{
			name: "empty",
			str:  "",
			want: nid.NID{},
		},
		{
			name:      "invalid",
			str:       "!invalid",
			wantPanic: true,
		},
		{
			name:      "no_name",
			str:       "000034o1ibe7u02570ak9evj9s",
			wantPanic: true,
		},
		{
			name:      "no_name_with_underscore",
			str:       "_000034o1ibe7u02570ak9evj9s",
			wantPanic: true,
		},
		{
			name:      "invalid_name",
			str:       "InvalidName_000034o1ibe7u02570ak9evj9s",
			wantPanic: true,
		},
		{
			name:      "invalid_base",
			str:       "book_!00034o1ibe7u02570ak9evj9s",
			wantPanic: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tc.wantPanic == (r == nil) {
					t.Errorf("MustParse() panic = %v; wantPanic = %v", r, tc.wantPanic)
				}
			}()

			id := nid.MustParse(tc.str)
			if id != tc.want {
				t.Errorf("MustParse() = %v; want = %v", id, tc.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tt := []struct { //nolint:dupl
		name    string
		str     string
		want    nid.NID
		wantErr bool
	}{
		{
			name: "valid",
			str:  "book_000034o1ibe7u02570ak9evj9s",
			want: nid.MustNaming("book").Apply(
				nid.Base{
					0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
					0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
				},
			),
		},
		{
			name: "zeros",
			str:  "book_00000000000000000000000000",
			want: nid.MustNaming("book").Apply(
				nid.Base{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
			),
		},
		{
			name: "empty",
			str:  "",
			want: nid.NID{},
		},
		{
			name:    "invalid",
			str:     "!invalid",
			wantErr: true,
		},
		{
			name:    "no_name",
			str:     "000034o1ibe7u02570ak9evj9s",
			wantErr: true,
		},
		{
			name:    "no_name_with_underscore",
			str:     "_000034o1ibe7u02570ak9evj9s",
			wantErr: true,
		},
		{
			name:    "invalid_name",
			str:     "InvalidName_000034o1ibe7u02570ak9evj9s",
			wantErr: true,
		},
		{
			name:    "invalid_base",
			str:     "book_!00034o1ibe7u02570ak9evj9s",
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			id, err := nid.Parse(tc.str)
			if tc.wantErr == (err == nil) {
				t.Errorf("Parse() panic = %v; wantErr = %v", err, tc.wantErr)
			}

			if id != tc.want {
				t.Errorf("Parse() = %v; want = %v", id, tc.want)
			}
		})
	}
}

func TestNIDString(t *testing.T) {
	tt := []struct {
		name string
		id   nid.NID
		want string
	}{
		{
			name: "empty",
			id:   nid.NID{},
			want: "",
		},
		{
			name: "empty_base",
			id:   nid.MustNaming("book").Apply(nid.Base{}),
			want: "",
		},
		{
			name: "normal",
			id: nid.MustNaming("book").Apply(
				nid.Base{
					0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
					0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
				},
			),
			want: "book_000034o1ibe7u02570ak9evj9s",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.id.String()
			if got != tc.want {
				t.Errorf("NID.String() = %s; want = %s", got, tc.want)
			}
		})
	}
}

func TestNIDMarshalText(t *testing.T) {
	tt := []struct {
		name string
		id   nid.NID
		want []byte
	}{
		{
			name: "empty",
			id:   nid.NID{},
			want: nil,
		},
		{
			name: "empty_base",
			id:   nid.MustNaming("book").Apply(nid.Base{}),
			want: nil,
		},
		{
			name: "normal",
			id: nid.MustNaming("book").Apply(
				nid.Base{
					0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
					0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
				},
			),
			want: []byte("book_000034o1ibe7u02570ak9evj9s"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.id.MarshalText()
			if err != nil {
				t.Errorf("NID.MarshalText() unexpected err = %v", err)
			}

			if !bytes.Equal(got, tc.want) {
				t.Errorf("NID.MarshalText() = %s; want = %s", got, tc.want)
			}
		})
	}
}

func TestNIDUnmarshalText(t *testing.T) {
	tt := []struct {
		name    string
		src     []byte
		want    nid.NID
		wantErr bool
	}{
		{
			name: "valid",
			src:  []byte("book_000034o1ibe7u02570ak9evj9s"),
			want: nid.MustNaming("book").Apply(
				nid.Base{
					0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
					0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
				},
			),
		},
		{
			name: "zeros",
			src:  []byte("book_00000000000000000000000000"),
			want: nid.MustNaming("book").Apply(
				nid.Base{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
			),
		},
		{
			name: "nil",
			src:  nil,
			want: nid.NID{},
		},
		{
			name: "empty",
			src:  []byte{},
			want: nid.NID{},
		},
		{
			name:    "invalid",
			src:     []byte("!invalid"),
			wantErr: true,
		},
		{
			name:    "no_name",
			src:     []byte("000034o1ibe7u02570ak9evj9s"),
			wantErr: true,
		},
		{
			name:    "no_name_with_underscore",
			src:     []byte("_000034o1ibe7u02570ak9evj9s"),
			wantErr: true,
		},
		{
			name:    "invalid_name",
			src:     []byte("InvalidName_000034o1ibe7u02570ak9evj9s"),
			wantErr: true,
		},
		{
			name:    "invalid_base",
			src:     []byte("book_!00034o1ibe7u02570ak9evj9s"),
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var id nid.NID

			err := id.UnmarshalText(tc.src)
			if tc.wantErr == (err == nil) {
				t.Errorf("NID.UnmarshalText() = %v; wantErr = %v", err, tc.wantErr)
			}

			if id != tc.want {
				t.Errorf("NID.UnmarshalText() = %v; want = %v", id, tc.want)
			}
		})
	}
}

func TestNIDMarshalJSON(t *testing.T) {
	tt := []struct {
		name string
		id   nid.NID
		want []byte
	}{
		{
			name: "empty",
			id:   nid.NID{},
			want: []byte("null"),
		},
		{
			name: "empty_base",
			id:   nid.MustNaming("book").Apply(nid.Base{}),
			want: []byte("null"),
		},
		{
			name: "normal",
			id: nid.MustNaming("book").Apply(
				nid.Base{
					0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
					0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
				},
			),
			want: []byte("\"book_000034o1ibe7u02570ak9evj9s\""),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.id.MarshalJSON()
			if err != nil {
				t.Errorf("NID.MarshalJSON() unexpected err = %v", err)
			}

			if !bytes.Equal(got, tc.want) {
				t.Errorf("NID.MarshalJSON() = %v; want = %v", got, tc.want)
			}
		})
	}
}

func TestNIDUnmarshalJSON(t *testing.T) {
	tt := []struct {
		name    string
		src     []byte
		want    nid.NID
		wantErr bool
	}{
		{
			name: "null",
			src:  []byte("null"),
			want: nid.NID{},
		},
		{
			name: "string",
			src:  []byte("\"book_000034o1ibe7u02570ak9evj9s\""),
			want: nid.MustNaming("book").Apply(
				nid.Base{
					0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
					0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
				},
			),
		},
		{
			name: "string_empty",
			src:  []byte("\"\""),
			want: nid.NID{},
		},
		{
			name:    "string_invalid_name",
			src:     []byte("\"Book_000034o1ibe7u02570ak9evj9s\""),
			wantErr: true,
		},
		{
			name:    "string_invalid_base",
			src:     []byte("\"book_0034o1ibe7u02570ak9evj9s\""),
			wantErr: true,
		},
		{
			name:    "number",
			src:     []byte("1.23"),
			wantErr: true,
		},
		{
			name:    "bool",
			src:     []byte("true"),
			wantErr: true,
		},
		{
			name:    "obj",
			src:     []byte("{}"),
			wantErr: true,
		},
		{
			name:    "array",
			src:     []byte("[]"),
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var id nid.NID

			err := id.UnmarshalJSON(tc.src)
			if tc.wantErr == (err == nil) {
				t.Errorf("NID.UnmarshalJSON() = %v; wantErr = %v", err, tc.wantErr)
			}

			if id != tc.want {
				t.Errorf("NID.UnmarshalJSON() = %v; want = %v", id, tc.want)
			}
		})
	}
}

func TestNIDValue(t *testing.T) {
	tt := []struct {
		name string
		id   nid.NID
		want driver.Value
	}{
		{
			name: "empty",
			id:   nid.NID{},
			want: nil,
		},
		{
			name: "valid",
			id: nid.MustNaming("book").Apply(
				nid.Base{
					0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
					0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
				},
			),
			want: "book_000034o1ibe7u02570ak9evj9s",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			value, err := tc.id.Value()
			if err != nil {
				t.Errorf("NID.Value() unexpected err = %v", err)
			}

			if !reflect.DeepEqual(value, tc.want) {
				t.Errorf("NID.Value() = %v; want = %v", value, tc.want)
			}
		})
	}
}

func TestNIDScan(t *testing.T) {
	tt := []struct {
		name    string
		src     any
		want    nid.NID
		wantErr bool
	}{
		{
			name: "nil",
			src:  nil,
			want: nid.NID{},
		},
		{
			name: "empty_bytes",
			src:  []byte{},
			want: nid.NID{},
		},
		{
			name: "nil_bytes",
			src:  []byte(nil),
			want: nid.NID{},
		},
		{
			name: "string",
			src:  "book_000034o1ibe7u02570ak9evj9s",
			want: nid.MustNaming("book").Apply(
				nid.Base{
					0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
					0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
				},
			),
		},
		{
			name: "bytes",
			src:  []byte("book_000034o1ibe7u02570ak9evj9s"),
			want: nid.MustNaming("book").Apply(
				nid.Base{
					0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
					0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
				},
			),
		},
		{
			name:    "string_invalid",
			src:     "book_0034o1ibe7u02570ak9evj9s",
			wantErr: true,
		},
		{
			name:    "bytes_invalid",
			src:     []byte("book_0034o1ibe7u02570ak9evj9s"),
			wantErr: true,
		},
		{
			name:    "int64",
			src:     int64(123),
			wantErr: true,
		},
		{
			name:    "float64",
			src:     float64(1.23),
			wantErr: true,
		},
		{
			name:    "bool",
			src:     true,
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var id nid.NID

			err := id.Scan(tc.src)
			if tc.wantErr == (err == nil) {
				t.Errorf("NID.Scan() = %v; wantErr = %v", err, tc.wantErr)
			}

			if id != tc.want {
				t.Errorf("NID.Scan() = %v; want = %v", id, tc.want)
			}
		})
	}
}
