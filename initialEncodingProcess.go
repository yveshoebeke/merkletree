package merkletree

func (ms *MerkleServer) initialEncodingProcess() {
	var (
		leaves [][]byte = *ms.leaves
	)

	for i := range leaves {
		leaves[i] = ms.hashGenerator(leaves[i])
	}
}
