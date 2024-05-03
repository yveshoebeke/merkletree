package merkletree

func (ms *MerkleServer) processDuplicateAndAppendRequest() error {
	var (
		started bool = false
	)

	for {
		// One remaining: Exit the loop. Merkle tree root determined.
		// Note: if initial data set is only 1 element, will continue
		//	so as to adhere to this Merkle tree discipline.
		if started && len(ms.Leaves) == 1 {
			break
		}
		started = true
		//  Adjust for odd number of leaves by duplicating last leave and appending it.
		//	- ie:
		//		[1] [2] [3] [4] [5] => [1] [2] [3] [4] [5] [5]
		if len(ms.Leaves)%2 == 1 {
			ms.Leaves = append(ms.Leaves, ms.Leaves[len(ms.Leaves)-1])
		}
		// Create combined (concatenated) hash of left and right (in couple),
		//	transform with requested algorithm (hash), and
		// 	store it in left and zero out right.
		//	- ie:
		// 		[1] [2] [3] [4] => [12] [0] [34] [0]
		for index := 0; index < len(ms.Leaves); index += 2 {
			ms.Leaves[index] = ms.hashGenerator(append(ms.Leaves[index][:], ms.Leaves[index+1][:]...))
			ms.Leaves[index+1] = []byte{}
		}

		// Remove 'nill' bytes.
		ms.removeNillBytes()
	}

	ms.ProcessResult = ms.Leaves[0]

	return nil
}