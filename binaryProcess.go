package merkletree

import (
	"math"
)

func (ms *MerkleServer) processBinaryRequest() error {
	var (
		leaves [][]byte = *ms.leaves
		index  int
	)

	// This could be removed.. TBD
	if ms.InitWithEncoding {
		for i := range leaves {
			leaves[i] = ms.hashGenerator(leaves[i])
		}
	}

	index = int(math.Pow(2, math.Ceil(math.Log2(float64(len(leaves))))))
	for ; index < len(leaves); index += 2 {
		leaves[index] = ms.hashGenerator(append(leaves[index][:], leaves[index+1][:]...))
		leaves[index+1] = []byte{}
	}

	leaves = removeNillBytes(leaves)

	for len(leaves) > 1 {
		for index = 0; index < len(leaves); index += 2 {
			leaves[index] = ms.hashGenerator(append(leaves[index][:], leaves[index+1][:]...))
			leaves[index+1] = []byte{}
		}

		leaves = removeNillBytes(leaves)
	}

	ms.ProcessResult = leaves[0]

	return nil
}
