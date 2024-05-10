package merkletree

import (
	"math"
)

func (ms *MerkleServer) processBinaryTreeRequest() error {
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
		// fmt.Println("[", len(ms.Leaves), "]")
		for index = 0; index < len(ms.Leaves); index += 2 {
			// - combine (concatenate) hash of left and right (in couple)
			// - encode it with requested algorithm
			// - Zero (nil) out the right element's value
			ms.Leaves[index] = ms.hashGenerator(append(ms.Leaves[index][:], ms.Leaves[index+1][:]...))
			ms.Leaves[index+1] = []byte{}
		}

		// Removenill bytes
		// fmt.Println("==>", len(ms.Leaves))

		ms.removeNillBytes(NopProcess, 0)
		// fmt.Println("===>", len(ms.Leaves))
	}

	ms.ProcessResult = ms.Leaves[0]

	return nil
}
