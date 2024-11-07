package nid

import "fmt"

var (
	ErrFailedParse = fmt.Errorf("nid: failed to parse")
	ErrInvalidName = fmt.Errorf("nid: invalid name")
)
