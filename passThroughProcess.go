package merkletree

func (ms *MerkleServer) processPassThroughRequest() error {
	var (
		leaves [][]byte
	)

	leaves = *ms.Leaves

	for len(leaves) > 1 {
		for index := 0; index < len(leaves); index += 2 {
			// if index to adjacent would overflow stop and leave last element alone,
			//	ie.: pass it through to next iteration.
			if index+1 >= len(leaves) {
				break
			}
			leaves[index] = ms.hashGenerator(append(leaves[index][:], leaves[index+1][:]...))
			leaves[index+1] = []byte{}
		}

		leaves = removeNillBytes(leaves)
	}

	ms.ProcessResult = leaves[0]

	return nil
}

// func (ms *MerkleServer) printHashes(msg string) {
// 	fmt.Printf("%s\n", msg)
// 	// fmt.Println(ms.Leaves)
// 	for _, l := range ms.Leaves {
// 		println(hex.EncodeToString(l))
// 	}
// 	fmt.Println("---")
// }
