package merkletree

import (
	"context"
	"encoding/json"
	"math"
	"sort"
)

func (ms *MerkleService) processBinaryTreeRequest(ctx context.Context) error {
	const ThisProcess = 2
	contextProcessType := ctx.Value(contextKeyRequestID)
	if contextProcessType != processTypes[ThisProcess] {
		return &InvalidContextProcessTypeErr{contextProcessType.(string)}
	}

	var index, startIndex int
	// Get starting index: if 2^x > length then idx = 2^x - length (See documentation for more details)
	startIndex = int(math.Pow(2, math.Ceil(math.Log2(float64(len(ms.Leaves)))))) - len(ms.Leaves)

	for index = startIndex; index < len(ms.Leaves); index += 2 {
		// - combine (concatenate) hash of left and right (in couple)
		// - encode it with requested algorithm
		// - Zero (nil) out the right element's value
		ms.Leaves[index] = ms.hashGenerator(append(ms.Leaves[index][:], ms.Leaves[index+1][:]...))
		ms.Leaves[index+1] = []byte{}
	}

	ms.removeNillBytes(BinaryTree, startIndex)

	for len(ms.Leaves) > 1 {
		for index = 0; index < len(ms.Leaves); index += 2 {
			// - combine (concatenate) hash of left and right (in couple)
			// - encode it with requested algorithm
			// - Zero (nil) out the right element's value
			ms.Leaves[index] = ms.hashGenerator(append(ms.Leaves[index][:], ms.Leaves[index+1][:]...))
			ms.Leaves[index+1] = []byte{}
		}

		// Removenill bytes
		ms.removeNillBytes(NopProcess, 0)
	}

	ms.ProcessResult = ms.Leaves[0]

	return nil
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
