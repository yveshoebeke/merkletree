package merkletree

// - Package: merkletree
//
// - Author:
//	 yves.hoebeke@bytesupply.com
//	 https://bytesupply.com/staff
//
// - Date: January 2024 (Haines City, FL USA)
//
// - Test, Benchmark: `go test -v; go test -branch=.`
//		- you can employ custom -detail flag in test, ex: `go test -detail`
//
// functions etc:
//
//	- DeriveRoot:
//	  --> This is the entry point to use this service.
//	  Merkle service configuration setup:
//		- Check if initial data (leaves) is available and of correct type.
//		- Check if requested algorithm (hash) exists and set related function.
//		- Check initial encoding flag
//		- Get initial array of hashes ("leaves")
//
//  - ProcessRequest:
//    Derive the root of a Merkle tree starting with data (i.e. transactions hashes in a block) as leaves.
//    Note: To enhance speed input data is transformed to a linked list.
////
//	- CurrentAlgorithmUsed:
//	  Returns the hash algorithm used.
//
//	- AvailableAlgorithms (cryptofuncs.go):
//	  Returns the hash algoritms available in this module.
//
//  - Process schematic:
//	  Note: Adjusting for odd number of leaves/branches by appending duplicate of last one.
//
// stop ------------> [12345555]          -> Merkle tree root (value out).
//                    /        \
//                [1234]       [5555]
//                /    \       /   \
//             [12]   [34]   [55] [55]    -> Branches.
//             [12]   [34]   [55]  ^
//             /  \   /  \   /  \
//            [1][2] [3][4] [5][5]
// start ---> [1][2] [3][4] [5] ^         -> Leaves (array data input).
//
// ^ = Duplicate hash to make leaf count even.
// A single leaf input will be processed as an odd number of leaves, ie: appended to itself.
//

import (
	"strings"
)

// Merkle tree object
type MerkleServer struct {
	leaves           *[][]byte  `json:"-"`
	HashTypeID       string     `json:"hashtype"`
	hashGenerator    CryptoFunc `json:"-"`
	InitWithEncoding bool       `json:"initWithEncoding"`
	ProofRequest     bool       `json:"proofrequest"`
	ProcessResult    []byte     `json:"root"`
	ProofResult      []byte     `json:"proof"`
}

/*
Entry Point

     ^ ^
    (O o)
-o00( . )00o-
-------------
*/
//   - Merkletree service configuration setup and start of request.
//
//     @params: name of algorithm to be used (string), initial data slice ([][]byte)
func GetRoot(algo string, data *[][]byte) {
	ms.DeriveRoot(algo, data)
}

func (ms *MerkleServer) DeriveRoot(algorithmRequested string, hashes *[][]byte, processFlags ...bool) ([]byte, error) {
	// Check if we got something to work with.
	if len(*hashes) == 0 {
		return []byte{}, &EmptyHashErr{}
	}

	// Check if requested algorithm is available.
	algorithmRequested = strings.ToUpper(algorithmRequested)
	if _, ok := AlgorithmRegistry[algorithmRequested]; !ok {
		return []byte{}, &UnknownAlgorithmErr{algorithmRequested}
	}

	// Set process flag
	initWithEncoding := If(len(processFlags) > 0, processFlags[0], false)

	// Initialize Merkel pertinents.
	ms = &MerkleServer{
		leaves:           hashes,
		HashTypeID:       algorithmRequested,
		hashGenerator:    AlgorithmRegistry[algorithmRequested],
		InitWithEncoding: initWithEncoding,
	}

	// Start process.
	if err := ms.processRequest(); err != nil {
		return []byte{}, err
	}

	// Return Merkle root from sliceProcessRequest
	return ms.ProcessResult, nil
}
