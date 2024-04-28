package merkletree

import (
	"slices"
)

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

// Remove elements with '0' byte content and collapse slice.
//   - ie:
//     [12] [0] [34] [0] => [12] [34]
func removeNillBytes(leaves [][]byte) [][]byte {
	return slices.DeleteFunc(leaves, func(leaf []byte) bool {
		return len(leaf) == 0
	})
}
