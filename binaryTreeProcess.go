package merkletree

import (
	"math"
)

func (ms *MerkleServer) processBinaryTreeRequest() error {
	var (
		index  int
		leaves [][]byte
	)

	leaves = *ms.Leaves

	// Get starting index: if 2^x > length -> idx = 2^x - length (See documentation for more details)
	for index = int(math.Pow(2, math.Ceil(math.Log2(float64(len(leaves)))))) - len(leaves); index < len(leaves); index += 2 {
		leaves[index] = ms.hashGenerator(append(leaves[index][:], leaves[index+1][:]...))
		leaves[index+1] = []byte{}
	}

	// Removenill bytes
	leaves = removeNillBytes(leaves)

	for len(leaves) > 1 {

		for index = 0; index < len(leaves); index += 2 {
			leaves[index] = ms.hashGenerator(append(leaves[index][:], leaves[index+1][:]...))
			leaves[index+1] = []byte{}
		}

		// Removenill bytes
		leaves = removeNillBytes(leaves)
	}

	ms.ProcessResult = leaves[0]
	return nil
}
