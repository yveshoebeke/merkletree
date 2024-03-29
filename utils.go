package merkletree

import "fmt"

// Custom error definitions.
// - empty data input
type EmptyHashErr struct{}

func (empty *EmptyHashErr) Error() string {
	return "empty hash argument"
}

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

type ProofErr struct{}

func (proof *ProofErr) Error() string {
	return "proof result did not match process result"
}

// Ternary operator
func If[T any](cond bool, trueValue, falseValue T) T {
	if cond {
		return trueValue
	}
	return falseValue
}

// Return the hash algorithm in use.
func (ms *MerkleServer) CurrentAlgorithmUsed() string {
	return ms.HashTypeID
}
