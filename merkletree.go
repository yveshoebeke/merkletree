package merkletree

//
// - Author:
//	 yves.hoebeke@bytesupply.com
//	 https://bytesupply.com/staff
//
// - Date: January 2024
//
// - Test, Benchmark: `go test -v; go test -branch=.`
//		- you can specify custom -detail flag in test, ex: `go test -detail`
//
// Functions:
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
// Helper/auxilary functions:
//
//	- CurrentAlgorithmUsed:
//		Returns the hash algorithm used.
//
//	- removeNillBytes:
//		Removes all empty []byte{} elements from hash slice.
//

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"
)

type contextKey int

const (
	contextKeyRequestID contextKey = iota
)

// GO routine response
type Response struct {
	err error
}
type processTypeFunction func(context.Context) error

// Merkle tree object
type MerkleServer struct {
	Leaves              [][]byte                    `json:"-"`
	HashTypeID          string                      `json:"hashtype"`
	hashGenerator       CryptoFunc                  `json:"-"`
	ProcessType         int                         `json:"processtype"`
	ProcessTypeRegistry map[int]processTypeFunction `json:"-"`
	InitWithEncoding    bool                        `json:"initWithEncoding"`
	ProofRequest        bool                        `json:"proofrequest"`
	ProcessResult       []byte                      `json:"root"`
	ProofResult         []byte                      `json:"proofresult"`
}

/*
Entry Point

- Merkletree service configuration setup and start of request.
*/
func DeriveRoot(algorithmRequested string, hashes [][]byte, processType int, initEncodingFlags bool) ([]byte, error) {

	// func DeriveRoot(algo string, data [][]byte, processType int, initEncode ...bool) ([]byte, error) {
	// ms := &MerkleServer{}
	// 	root, err := ms.GetMerkletreeRoot(algo, data, processType, initEncode)
	// 	if err != nil {
	// 		return []byte{}, err
	// 	}

	// 	return root, nil
	// }

	// func (ms *MerkleServer) GetMerkletreeRoot(algorithmRequested string, hashes [][]byte, processType int, initEncodingFlags ...bool) ([]byte, error) {
	// type processTypeFunction func(context.Context) error

	// Check if we got something to work with.
	if len(hashes) == 0 {
		return []byte{}, &EmptyListErr{}
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
	//initWithEncoding := If(len(initEncodingFlags) > 0, initEncodingFlags[0], []bool{false}[0])

	// Initialize merkle pertinents.
	ms := &MerkleServer{
		Leaves:              hashes,
		HashTypeID:          algorithmRequested,
		hashGenerator:       AlgorithmRegistry[algorithmRequested],
		ProcessType:         processType,
		ProcessTypeRegistry: make(map[int]processTypeFunction),
		InitWithEncoding:    initEncodingFlags,
	}

	// Process type registry
	ms.ProcessTypeRegistry[0] = ms.processPassThroughRequest
	ms.ProcessTypeRegistry[1] = ms.processDuplicateAndAppendRequest
	ms.ProcessTypeRegistry[2] = ms.processBinaryTreeRequest

	// Hash first branch (input hash slice) if requested.
	if ms.InitWithEncoding {
		for i := range ms.Leaves {
			ms.Leaves[i] = ms.hashGenerator(ms.Leaves[i])
		}
	}

	// Set context process id and timeout criteria
	ctx := context.WithValue(context.Background(), contextKeyRequestID, processTypes[processType])
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*ProcessTimeoutMilliSecs)
	defer cancel()

	// Response channel
	resch := make(chan Response)

	// Execute desired processtype
	go func() {
		err := ms.ProcessTypeRegistry[processType](ctx)
		resch <- Response{err: err}
	}()

	// Get results, errors from response channel
	select {
	case <-ctx.Done():
		return []byte{}, fmt.Errorf("timed out: %+w", ctx.Err())
	case resp := <-resch:
		if resp.err != nil {
			return []byte{}, resp.err
		}
	}

	return ms.ProcessResult, nil
}

// Return the hash algorithm in use.
func (ms *MerkleServer) CurrentAlgorithmUsed() string {
	return ms.HashTypeID
}

// Remove elements with nill byte content and collapse slice.
//   - ie:
//     [12] [0] [34] [0] => [12] [34]
func (ms *MerkleServer) removeNillBytes(processType int, startValue ...int) {
	zeroStart := []int{0}
	start := If(len(startValue) > 0, startValue[0], zeroStart[0])
	switch processType {
	case BinaryTree:
		l := len(ms.Leaves)
		found := true
		for found {
			found = false
			for i := start; i < l; i++ {
				if len(ms.Leaves[i]) == 0 {
					ms.Leaves = append(ms.Leaves[:i], ms.Leaves[i+1:]...)
					found = true
					l--
				}
			}
		}

	default:
		ms.Leaves = slices.DeleteFunc(ms.Leaves, func(leaf []byte) bool {
			return len(leaf) == 0
		})
	}
}
