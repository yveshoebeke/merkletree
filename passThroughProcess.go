package merkletree

func (ms *MerkleServer) processPassThroughRequest() error {
	var (
		leaves [][]byte = *ms.leaves
	)

	for len(leaves) > 1 {
		for index := 0; index < len(leaves); index += 1 {
			leaves[index] = ms.hashGenerator(append(leaves[index][:], leaves[index+1][:]...))
			leaves[index+1] = []byte{}

			index++
			leaves = removeNillBytes(leaves)
		}
	}

	ms.ProcessResult = leaves[0]

	return nil
}
