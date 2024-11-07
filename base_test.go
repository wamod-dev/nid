package nid_test

import (
	"bytes"
	"crypto/rand"
	"database/sql/driver"
	"os"
	"reflect"
	"testing"
	"time"

	"go.wamod.dev/nid"
)

func TestMustParseBase(t *testing.T) {
	tt := []struct {
		name      string
		str       string
		want      nid.Base
		wantPanic bool
	}{
		{
			name:      "panic",
			str:       "!invalid",
			wantPanic: true,
		},
		{
			name: "valid",
			str:  "000034o1ibe7u02570ak9evj9s",
			want: nid.Base{
				0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
				0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
			},
		},
		{
			name: "empty",
			str:  "",
			want: nid.Base{},
		},
		{
			name: "zero",
			str:  "00000000000000000000000000",
			want: nid.Base{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tc.wantPanic == (r == nil) {
					t.Errorf("MustParseBase() panic = %v; wantPanic = %v", r, tc.wantPanic)
				}
			}()

			got := nid.MustParseBase(tc.str)
			if got != tc.want {
				t.Errorf("MustParseBase() = %v; want = %v", got, tc.want)
			}
		})
	}
}

type failReader struct {
	err error
}

func (r failReader) Read(_ []byte) (int, error) {
	return 0, r.err
}

func TestNewBaseAt(t *testing.T) {
	randReader := rand.Reader

	tt := []struct {
		name      string
		ts        time.Time
		wantTime  time.Time
		wantPanic bool

		before func()
		after  func()
	}{
		{
			name:     "normal",
			ts:       time.UnixMilli(12345),
			wantTime: time.UnixMilli(12345),
		},
		{
			name:      "bad_rand_reader",
			ts:        time.UnixMilli(12345),
			wantPanic: true,
			before: func() {
				rand.Reader = failReader{os.ErrClosed}
			},
			after: func() {
				rand.Reader = randReader
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.before != nil {
				tc.before()
			}

			defer func() {
				r := recover()
				if tc.wantPanic == (r == nil) {
					t.Errorf("MustParseBase() panic = %v; wantPanic = %v", r, tc.wantPanic)
				}

				if tc.after != nil {
					tc.after()
				}
			}()

			base := nid.NewBaseAt(tc.ts)

			if base.Time() != tc.wantTime {
				t.Errorf("NewBaseAt() time = %v; wantTime = %v", base.Time(), tc.wantTime)
			}
		})
	}
}

func TestBaseValue(t *testing.T) {
	tt := []struct {
		name string
		base nid.Base
		want driver.Value
	}{
		{
			name: "empty",
			base: nid.Base{},
			want: nil,
		},
		{
			name: "zeroes",
			base: nid.Base{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			want: nil,
		},
		{
			name: "normal",
			base: nid.Base{
				0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
				0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
			},
			want: []byte{
				0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
				0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			value, err := tc.base.Value()
			if err != nil {
				t.Fatalf("Base.Value() unexpected err = %v", err)
			}

			if !reflect.DeepEqual(value, tc.want) {
				t.Errorf("Base.Value() = %v; want = %v", value, tc.want)
			}
		})
	}
}

func TestBaseScan(t *testing.T) {
	tt := []struct {
		name    string
		src     any
		want    nid.Base
		wantErr bool
	}{
		{
			name: "nil",
			src:  nil,
			want: nid.Base{},
		},
		{
			name: "string",
			src:  "000034o1ibe7u02570ak9evj9s",
			want: nid.Base{
				0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
				0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
			},
		},
		{
			name: "string_zeros",
			src:  "00000000000000000000000000",
			want: nid.Base{},
		},
		{
			name: "string_empty",
			src:  "",
			want: nid.Base{},
		},
		{
			name: "bytes",
			src: []byte{
				0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
				0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
			},
			want: nid.Base{
				0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
				0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
			},
		},
		{
			name:    "bytes_invalid_len",
			src:     []byte{0x0, 0x0, 0x1, 0x93, 0x1, 0x92, 0xdc, 0x7f},
			wantErr: true,
		},
		{
			name: "bytes_zeros",
			src: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			want: nid.Base{},
		},
		{
			name: "bytes_empty",
			src:  []byte{},
			want: nid.Base{},
		},
		{
			name: "bytes_nil",
			src:  []byte(nil),
			want: nid.Base{},
		},
		{
			name:    "int64",
			src:     int64(123),
			wantErr: true,
		},
		{
			name:    "float64",
			src:     float64(123),
			wantErr: true,
		},
		{
			name:    "bool",
			src:     true,
			wantErr: true,
		},
		{
			name:    "time",
			src:     time.UnixMilli(123),
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var base nid.Base

			err := base.Scan(tc.src)
			if tc.wantErr == (err == nil) {
				t.Errorf("Base.Scan() = %v; wantErr = %v", err, tc.wantErr)
			}

			if base != tc.want {
				t.Errorf("Base.Scan() = %v; want = %v", base, tc.want)
			}
		})
	}
}

func TestBaseMarshalText(t *testing.T) {
	base := nid.Base{
		0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
		0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
	}
	want := []byte("000034o1ibe7u02570ak9evj9s")

	got, err := base.MarshalText()
	if err != nil {
		t.Fatalf("Base.MarshalText() unexpected err = %v", err)
	}

	if !bytes.Equal(got, want) {
		t.Errorf("Base.MarshalText() = %v; want = %v", got, want)
	}
}

func TestBaseUnmarshalText(t *testing.T) {
	tt := []struct {
		name    string
		src     []byte
		want    nid.Base
		wantErr bool
	}{
		{
			name: "nil",
			src:  nil,
			want: nid.Base{},
		},
		{
			name: "empty",
			src:  []byte{},
			want: nid.Base{},
		},
		{
			name: "valid",
			src:  []byte("000034o1ibe7u02570ak9evj9s"),
			want: nid.Base{
				0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
				0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
			},
		},
		{
			name:    "invalid_len",
			src:     []byte("000034o1ibe7u02570ak9evj"),
			wantErr: true,
		},
		{
			name:    "invalid_encoding",
			src:     []byte("_00034o1ibe7u02570ak9evj9j"),
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var base nid.Base

			err := base.UnmarshalText(tc.src)
			if tc.wantErr == (err == nil) {
				t.Errorf("Base.UnmarshalText() = %v; wantErr = %v", err, tc.wantErr)
			}

			if base != tc.want {
				t.Errorf("Base.UnmarshalText() = %v; want = %v", base, tc.want)
			}
		})
	}
}

func TestBaseMarshalJSON(t *testing.T) {
	tt := []struct {
		name string
		base nid.Base
		want []byte
	}{
		{
			name: "empty",
			base: nid.Base{},
			want: []byte("null"),
		},
		{
			name: "valid",
			base: nid.Base{
				0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
				0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
			},
			want: []byte("\"000034o1ibe7u02570ak9evj9s\""),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.base.MarshalJSON()
			if err != nil {
				t.Errorf("Base.MarshalJSON() = unexpected error = %v", err)
			}

			if !bytes.Equal(got, tc.want) {
				t.Errorf("Base.MarshalJSON() = %v; want = %v", got, tc.want)
			}
		})
	}
}

func TestBaseUnmarshalJSON(t *testing.T) {
	tt := []struct {
		name    string
		src     []byte
		want    nid.Base
		wantErr bool
	}{
		{
			name: "null",
			src:  []byte("null"),
			want: nid.Base{},
		},
		{
			name: "valid",
			src:  []byte("\"000034o1ibe7u02570ak9evj9s\""),
			want: nid.Base{
				0x00, 0x00, 0x01, 0x93, 0x01, 0x92, 0xdc, 0x7f,
				0x00, 0x45, 0x38, 0x15, 0x44, 0xbb, 0xf3, 0x4f,
			},
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
			name:    "empty",
			src:     []byte{},
			wantErr: true,
		},
		{
			name:    "object",
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
			var base nid.Base

			err := base.UnmarshalJSON(tc.src)
			if tc.wantErr == (err == nil) {
				t.Errorf("Base.UnmarshalJSON() = %v; wantErr = %v", err, tc.wantErr)
			}

			if base != tc.want {
				t.Errorf("Base.UnmarshalJSON() = %v; want = %v", base, tc.want)
			}
		})
	}
}
