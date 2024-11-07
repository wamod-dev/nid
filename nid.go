// Package nid implements named unique sortable identifiers.
//
// It's designed to provide a way to work with unique identifiers that are more human-readable.
// The identifiers are sortable and can be used in the database.
// The package also provides a way to convert the identifiers to and from the JSON, Text, and database values.
package nid

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

// NID is a named identifier.
// It consists of the name and the base identifier that is sortable and unique.
type NID struct {
	name string
	base Base
}

// MustParse is a helper to parse named [NID]. It panics if given string is invalid.
func MustParse(str string) NID {
	id, err := Parse(str)
	if err != nil {
		panic(err)
	}

	return id
}

// Parse the named ID from the string.
func Parse(str string) (dst NID, err error) {
	err = dst.UnmarshalText([]byte(str))

	return
}

// Name returns the name of the ID.
func (id NID) Name() string {
	return id.name
}

// Base returns the base identifier.
func (id NID) Base() Base {
	return id.base
}

// String returns the string representation of the ID.
// The format is "<name>_<id>".
func (id NID) String() string {
	bytes, _ := id.MarshalText()

	return string(bytes)
}

// MarshalText returns the text representation of the ID.
func (id NID) MarshalText() ([]byte, error) {
	if len(id.name) == 0 || id.base.Empty() {
		return []byte{}, nil
	}

	str := strings.Join([]string{id.name, id.base.String()}, "_")

	return []byte(str), nil
}

// UnmarshalText parses the ID from the text.
func (id *NID) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		*id = NID{}

		return nil
	}

	str := string(data)

	cut := strings.LastIndex(str, "_")
	if cut <= 0 || cut == len(str)-1 {
		return fmt.Errorf("%w: invalid named identifier: %q", ErrFailedParse, str)
	}

	name := str[:cut]
	if !validateName(name) {
		return fmt.Errorf("%w: identifier name must be a non-empty snake_case string: %q", ErrFailedParse, name)
	}

	err := id.base.UnmarshalText([]byte(str[cut+1:]))
	if err != nil {
		return err
	}

	if id.base.Empty() {
		id.name = ""
	} else {
		id.name = name
	}

	return nil
}

// MarshalJSON returns the JSON representation of the ID.
func (id NID) MarshalJSON() ([]byte, error) {
	if id.Empty() {
		return []byte("null"), nil
	}

	return json.Marshal(id.String())
}

// UnmarshalJSON parses the ID from the JSON.
func (id *NID) UnmarshalJSON(src []byte) error {
	if bytes.Equal(src, []byte("null")) {
		*id = NID{}

		return nil
	}

	var str string

	err := json.Unmarshal(src, &str)
	if err != nil {
		return err
	}

	return id.UnmarshalText([]byte(str))
}

// Empty returns true if the ID is empty.
func (id NID) Empty() bool {
	return len(id.name) == 0 || id.base.Empty()
}

// Value returns the driver value.
func (id NID) Value() (driver.Value, error) {
	if id.Empty() {
		return nil, nil
	}

	return id.String(), nil
}

// Scan the value into the ID.
func (id *NID) Scan(src any) error {
	if src == nil {
		id.base = Base{}
		id.name = ""

		return nil
	}

	switch src := src.(type) {
	case string:
		return id.UnmarshalText([]byte(src))
	case []byte:
		return id.UnmarshalText(src)
	default:
		return fmt.Errorf("%w: invalid scan source: %T", ErrFailedParse, src)
	}
}
