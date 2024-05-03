package merkletree

import (
	"math"
)

func (ms *MerkleServer) processBinaryTreeRequest() error {
	// Get starting index: if 2^x > length -> idx = 2^x - length (See documentation for more details)
	index := int(math.Pow(2, math.Ceil(math.Log2(float64(len(ms.Leaves)))))) - len(ms.Leaves)
	for ; index < len(ms.Leaves); index += 2 {
		ms.Leaves[index] = ms.hashGenerator(append(ms.Leaves[index][:], ms.Leaves[index+1][:]...))
		ms.Leaves[index+1] = []byte{}
	}

	// Removenill bytes
	ms.removeNillBytes()

	// Continue - each branch is now adhering to binary tree model.
	for len(ms.Leaves) > 1 {
		for index = 0; index < len(ms.Leaves); index += 2 {
			ms.Leaves[index] = ms.hashGenerator(append(ms.Leaves[index][:], ms.Leaves[index+1][:]...))
			ms.Leaves[index+1] = []byte{}
		}

		ms.removeNillBytes()
	}

	ms.ProcessResult = ms.Leaves[0]

	return nil
}
