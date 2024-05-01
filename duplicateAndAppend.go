package merkletree

func (ms *MerkleServer) processDuplicateAndAppendRequest() error {
	var (
		leaves  [][]byte = *ms.leaves
		started bool     = false
		index   int
	)

	for {
		// One remaining: Exit loop. Merkle tree root made.
		// Note: if initial data set is only 1 element, will continue
		//	so as to adhere to the Merkle tree discipline.
		if started && len(leaves) == 1 {
			break
		}
		started = true
		//  Adjust for odd number of leaves.
		//	- ie:
		//		[1] [2] [3] => [1] [2] [3] [3]
		if len(leaves)%2 == 1 {
			leaves = append(leaves, leaves[len(leaves)-1])
		}
		// Create combined (concatenated) hash of left and right (in couple),
		//	transform with requested algorithm (hash), and
		// 	store it in left and zero out right.
		//	- ie:
		// 		[1] [2] [3] [4] => [12] [0] [34] [0]
		for index = 0; index < len(leaves); index += 2 {
			leaves[index] = ms.hashGenerator(append(leaves[index][:], leaves[index+1][:]...))
			leaves[index+1] = []byte{}
		}

		leaves = removeNillBytes(leaves)
	}

	ms.ProcessResult = leaves[0]

	return nil
}
