package merkletree

import (
	"fmt"
)

// Custom error definitions.

// - empty data input
type EmptyListErr struct{}

func (empty *EmptyListErr) Error() string {
	return "empty list"
}

// - unknown hash algorithm requested (not in crypto function registry)
type UnknownAlgorithmErr struct {
	algorithmRequest string
}

func (unknown *UnknownAlgorithmErr) Error() string {
	return fmt.Sprintf("unknown hash algorithm: %s", unknown.algorithmRequest)
}

// - invalid process type value argument
type InvalidProcessTypeErr struct{}

func (empty *InvalidProcessTypeErr) Error() string {
	return "invalid process type argument"
}
