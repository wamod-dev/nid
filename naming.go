package nid

import (
	"fmt"
	"time"
)

// Naming provides a way to create, update and validate the [NID]s.
type Naming struct {
	name string
}

// MustNaming is a helper to create Namer from the name. It panics if the name is invalid.
func MustNaming(name string) Naming {
	n, err := NewNaming(name)
	if err != nil {
		panic(err)
	}

	return n
}

// NewNaming creates a new [Naming] from the name. It returns an error if the name is invalid.
// The name must be a non-empty snake_case string, e.g. "user" or "user_profile".
func NewNaming(name string) (Naming, error) {
	if !validateName(name) {
		return Naming{}, fmt.Errorf("%w: must be a non-empty snake_case string: %s", ErrInvalidName, name)
	}

	return Naming{name}, nil
}

// initialized the [Naming] has a name.
func (n Naming) initialized() {
	if n.name == "" {
		panic("nid: identifier naming was not initialized")
	}
}

// New creates a new [NID] at the current time.
func (n Naming) New() NID {
	n.initialized()

	return NID{
		name: n.name,
		base: NewBase(),
	}
}

// NewAt creates a new [NID] at the given time.
func (n Naming) NewAt(ts time.Time) NID {
	n.initialized()

	return NID{
		name: n.name,
		base: NewBaseAt(ts),
	}
}

// Is checks if the name of the [NID] matches namer's name.
func (n Naming) Is(id NID) bool {
	n.initialized()

	return n.name == id.Name()
}

// Apply create a new [NID] from the given [Base].
func (n Naming) Apply(base Base) NID {
	n.initialized()

	if base.Empty() {
		return NID{}
	}

	return NID{
		name: n.name,
		base: base,
	}
}

// Update the name of the [NID] with the namer's name.
func (n Naming) Update(id NID) NID {
	n.initialized()

	if id.Empty() {
		return NID{}
	}

	return NID{
		name: n.name,
		base: id.base,
	}
}

func validateName(str string) bool {
	ok := false

	for i, r := range str {
		switch {
		case r == '_':
			if !ok {
				return false
			}

			ok = false
		case ('0' <= r && r <= '9'):
			if i == 0 {
				return false
			}

			ok = true
		case ('a' <= r && r <= 'z'):
			ok = true
		default:
			return false
		}
	}

	return ok
}
