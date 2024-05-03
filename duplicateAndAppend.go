package merkletree

func (ms *MerkleServer) processDuplicateAndAppendRequest() error {
	var (
		leaves  [][]byte
		started bool = false
	)

	leaves = *ms.Leaves

	for {
		// One remaining: Exit the loop. Merkle tree root determined.
		// Note: if initial data set is only 1 element, will continue
		//	so as to adhere to this Merkle tree discipline.
		if started && len(leaves) == 1 {
			break
		}
		started = true
		//  Adjust for odd number of leaves by duplicating last leave and appending it.
		//	- ie:
		//		[1] [2] [3] [4] [5] => [1] [2] [3] [4] [5] [5]
		if len(leaves)%2 == 1 {
			leaves = append(leaves, leaves[len(leaves)-1])
		}
		// Create combined (concatenated) hash of left and right (in couple),
		//	transform with requested algorithm (hash), and
		// 	store it in left and zero out right.
		//	- ie:
		// 		[1] [2] [3] [4] => [12] [0] [34] [0]
		for index := 0; index < len(leaves); index += 2 {
			leaves[index] = ms.hashGenerator(append(leaves[index][:], leaves[index+1][:]...))
			leaves[index+1] = []byte{}
		}

		// Remove 'nill' bytes.
		leaves = removeNillBytes(leaves)
	}

	ms.ProcessResult = leaves[0]

	return nil
}
