package merkletree

func (ms *MerkleServer) processPassThroughRequest() error {
	for len(ms.Leaves) > 1 {
		for index := 0; index < len(ms.Leaves); index += 2 {
			// - if index to adjacent would overflow stop and leave last element alone,
			//	wow: pass it through to next branch iteration.
			if index+1 >= len(ms.Leaves) {
				break
			}

			// - combine (concatenate) hash of left and right (in couple)
			// - encode it with requested algorithm
			// - Zero (nil) out the right element's value
			ms.Leaves[index] = ms.hashGenerator(append(ms.Leaves[index][:], ms.Leaves[index+1][:]...))
			ms.Leaves[index+1] = []byte{}
		}

		ms.removeNillBytes()
	}

	ms.ProcessResult = ms.Leaves[0]

	return nil
}
