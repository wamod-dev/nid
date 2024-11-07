package nid

import (
	"bytes"
	"crypto/rand"
	"database/sql/driver"
	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"
)

const (
	timeLen = 8
	randLen = 8
	baseLen = timeLen + randLen
	encStr  = "0123456789abcdefghijklmnopqrstuv"
)

var encoding = base32.NewEncoding(encStr).WithPadding(base32.NoPadding) //nolint:gochecknoglobals

// Base of the [NID] with a time and random part.
// The time part is 8 bytes and the random part is 8 bytes.
// The total length is 16 bytes. The [Base] is sortable and unique.
type Base [baseLen]byte

// NewBaseAt creates a new [Base] for the given time.
// The time is in nanoseconds.
func NewBaseAt(ts time.Time) Base {
	var dst Base

	binary.BigEndian.PutUint64(dst[:timeLen], uint64(ts.UnixMilli()))

	if _, err := rand.Read(dst[timeLen:]); err != nil {
		panic(err)
	}

	return dst
}

// NewBase creates a new [Base] at the current time.
func NewBase() Base {
	return NewBaseAt(time.Now())
}

// ParseBaseBytes parses the [Base] from the bytes.
func ParseBaseBytes(src []byte) (dst Base, err error) {
	err = dst.UnmarshalText(src)

	return
}

// ParseBase parses the [Base] from the string.
func ParseBase(src string) (Base, error) {
	return ParseBaseBytes([]byte(src))
}

// MustParseBase is a helper to parse [Base]. It panics if the base is invalid.
func MustParseBase(src string) Base {
	base, err := ParseBase(src)
	if err != nil {
		panic(err)
	}

	return base
}

// Unix returns the time of the [Base] in seconds.
func (base Base) UnixMilli() int64 {
	ts := binary.BigEndian.Uint64(base[:timeLen])

	return int64(ts) //nolint:gosec
}

// Time returns the time of the [Base].
func (base Base) Time() time.Time {
	return time.Unix(0, base.UnixMilli()*int64(time.Millisecond))
}

// Value returns the driver value of the [Base].
func (base Base) Value() (driver.Value, error) {
	if base.Empty() {
		return nil, nil
	}

	return base.Bytes(), nil
}

// Scan scans the value into the [Base].
func (base *Base) Scan(src any) (err error) {
	if src == nil {
		*base = Base{}

		return nil
	}

	switch src := src.(type) {
	case string:
		return base.UnmarshalText([]byte(src))
	case []byte:
		if l := len(src); l == 0 {
			*base = Base{}

			return nil
		} else if l != baseLen {
			return fmt.Errorf("%w: invalid base id length: %d", ErrFailedParse, len(src))
		}

		copy(base[:], src)

		return nil
	default:
		return fmt.Errorf("%w: invalid scan source: %T", ErrFailedParse, src)
	}
}

// MarshalText returns the text representation of the [Base].
func (base Base) MarshalText() ([]byte, error) {
	dst := make([]byte, encoding.EncodedLen(baseLen))
	encoding.Encode(dst, base[:])

	return dst, nil
}

// UnmarshalText parses the [Base] from the text.
func (base *Base) UnmarshalText(src []byte) error {
	l := len(src)
	if l == 0 {
		*base = Base{}

		return nil
	} else if l != encoding.EncodedLen(baseLen) {
		return fmt.Errorf("%w: invalid base id length: %d", ErrFailedParse, l)
	}

	var dst Base

	_, err := encoding.Decode(dst[:], src)
	if err != nil {
		return fmt.Errorf("%w: invalid base encoding: %w", ErrFailedParse, err)
	}

	*base = dst

	return nil
}

// MarshalJSON returns the JSON representation of the [Base].
func (base Base) MarshalJSON() ([]byte, error) {
	if base.Empty() {
		return []byte("null"), nil
	}

	return json.Marshal(base.String())
}

// UnmarshalJSON parses the [Base] from the JSON.
func (base *Base) UnmarshalJSON(src []byte) error {
	if bytes.Equal(src, []byte("null")) {
		*base = Base{}

		return nil
	}

	var str string

	err := json.Unmarshal(src, &str)
	if err != nil {
		return err
	}

	return base.UnmarshalText([]byte(str))
}

// Bytes returns the bytes of the [Base].
func (base Base) Bytes() []byte {
	return base[:]
}

// Empty returns true if the [Base] is empty.
func (base Base) Empty() bool {
	for _, b := range base {
		if b != 0 {
			return false
		}
	}

	return true
}

// String returns the string representation of the [Base].
func (base Base) String() string {
	bytes, _ := base.MarshalText()

	return string(bytes)
}
