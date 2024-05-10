package merkletree

//  - If:
//		General ternary operator.
//
//	- CurrentAlgorithmUsed:
//		Returns the hash algorithm used.
//
//	- removeNillBytes:
//		Removes all empty []byte{} elements from hash slice.
//
//	- AvailableAlgorithms (cryptofuncs.go):
//		Returns the hash algoritms available in this module.
//
import (
	"encoding/json"
	"slices"
	"sort"
)

const (
	NopProcess  = -1
	PassThrough = 0
	DupeAppend  = 1
	BinaryTree  = 2
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

// Remove elements with nill byte content and collapse slice.
//   - ie:
//     [12] [0] [34] [0] => [12] [34]
func (ms *MerkleServer) removeNillBytes(processType int, startValue ...int) {
	zeroStart := []int{0}
	start := If(len(startValue) > 0, startValue[0], zeroStart[0])
	switch processType {
	case BinaryTree:
		l := len(ms.Leaves)
		found := true
		for found {
			found = false
			for i := start; i < l; i++ {
				if len(ms.Leaves[i]) == 0 {
					ms.Leaves = append(ms.Leaves[:i], ms.Leaves[i+1:]...)
					found = true
					l--
				}
			}
		}

	default:
		ms.Leaves = slices.DeleteFunc(ms.Leaves, func(leaf []byte) bool {
			return len(leaf) == 0
		})
	}
}

// Returns json string with all available hash functions.
func AvailableAlgorithms() (string, error) {
	type availableJson struct {
		Algorithms []string `json:"algorithms"`
	}
	jsonResult := &availableJson{}

	for algorithmName := range AlgorithmRegistry {
		jsonResult.Algorithms = append(jsonResult.Algorithms, algorithmName)
	}

	sort.Strings(jsonResult.Algorithms)
	jsonEncodedAlgorithms, err := json.Marshal(jsonResult)
	if err != nil {
		return "", err
	}

	return string(jsonEncodedAlgorithms), nil
}
