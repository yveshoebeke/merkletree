package merkletree

func (ms *MerkleServer) processPassThroughRequest() error {
	for len(ms.Leaves) > 1 {
		for index := 0; index < len(ms.Leaves); index += 2 {
			// if index to adjacent would overflow stop and leave last element alone,
			//	ie.: pass it through to next iteration.
			if index+1 >= len(ms.Leaves) {
				break
			}
			ms.Leaves[index] = ms.hashGenerator(append(ms.Leaves[index][:], ms.Leaves[index+1][:]...))
			ms.Leaves[index+1] = []byte{}
		}

		ms.removeNillBytes()
	}

	ms.ProcessResult = ms.Leaves[0]

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
