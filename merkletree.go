package merkletree

// - Package: merkletree
//
// - Author:
//	 yves.hoebeke@bytesupply.com
//	 https://bytesupply.com/staff
//
// - Date: January 2024
//
// - Test, Benchmark: `go test -v; go test -branch=.`
//		- you can employ custom -detail flag in test, ex: `go test -detail`
//
// functions etc:
//
//	- DeriveRoot:
//	  --> This is the entry point to use this service.
//
//	- GetMerkletreeRoot:
//	  	Merkletree service configuration setup:
//			- Check if initial data (leaves) is available and of correct type.
//			- Validate requested algorithm (hash) exists and set related function.
//			- Check initial encoding flag.
//			- Validate process type requested.
//			- Init Merkletree Structure.
//			- Hash all elements of first branch, if requested.
//			- Execute requested tree manipulation request.
//
import (
	"strings"
)

// Merkle tree object
type MerkleServer struct {
	Leaves           [][]byte   `json:"-"`
	HashTypeID       string     `json:"hashtype"`
	hashGenerator    CryptoFunc `json:"-"`
	ProcessType      int        `json:"processtype"`
	InitWithEncoding bool       `json:"initWithEncoding"`
	ProofRequest     bool       `json:"proofrequest"`
	ProcessResult    []byte     `json:"root"`
	ProofResult      []byte     `json:"proofresult"`
}

/*
Entry Point

	 ^ ^
	(O o)

-o00( . )00o-
-------------

  - Merkletree service configuration setup and start of request.
*/
func DeriveRoot(algo string, data [][]byte, processType int, initEncode bool) ([]byte, error) {
	ms := &MerkleServer{}
	root, err := ms.GetMerkletreeRoot(algo, data, processType, initEncode)
	if err != nil {
		return []byte{}, err
	}

	return root, nil
}

func (ms *MerkleServer) GetMerkletreeRoot(algorithmRequested string, hashes [][]byte, processType int, initEncodingFlags bool) ([]byte, error) {
	// Check if we got something to work with.
	if len(hashes) == 0 {
		return []byte{}, &EmptyHashErr{}
	}

	// Check if requested algorithm is available.
	algorithmRequested = strings.ToUpper(algorithmRequested)
	if _, ok := AlgorithmRegistry[algorithmRequested]; !ok {
		return []byte{}, &UnknownAlgorithmErr{algorithmRequested}
	}

	// Set/check process type request.
	if processType < 0 || processType > 2 {
		return []byte{}, &InvalidProcessTypeErr{}
	}

	// Set process flag.
	// initWithEncoding := If(len(initEncodingFlags) > 0, initEncodingFlags[0], false)
	initWithEncoding := initEncodingFlags

	// Initialize Merkel pertinents.
	ms = &MerkleServer{
		Leaves:           hashes,
		HashTypeID:       algorithmRequested,
		hashGenerator:    AlgorithmRegistry[algorithmRequested],
		ProcessType:      processType,
		InitWithEncoding: initWithEncoding,
	}

	// Hash first branch (input hash slice) if requested.
	if ms.InitWithEncoding {
		for i := range ms.Leaves {
			ms.Leaves[i] = ms.hashGenerator(ms.Leaves[i])
		}
	}
	// Start requested process. Unbalanced trees will be handled according
	//	to the specific desired process logic.
	//	- 0: duplicate and append last hash element to current branch.
	//	- 1: pass last hash element of odd length branch to next.
	//	- 2: convert to binary tree.
	switch processType {
	case 0:
		if err := ms.processDuplicateAndAppendRequest(); err != nil {
			return []byte{}, err
		}
	case 1:
		if err := ms.processPassThroughRequest(); err != nil {
			return []byte{}, err
		}
	case 2:
		if err := ms.processBinaryTreeRequest(); err != nil {
			return []byte{}, err
		}

	}

	// Return Merkle root from process
	return ms.ProcessResult, nil
}
