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
	"slices"
	"strings"
	"time"
)

// General purpose app constants
const (
	NopProcess                         = -1
	PassThrough                        = 0
	DupeAppend                         = 1
	BinaryTree                         = 2
	ProcessTimeoutMilliSecs            = 100
	contextKeyRequestID     contextKey = iota
)

// CTX key values
var (
	processTypes = [3]string{"PAS-THRU", "DUP-APND", "BIN-TREE"}
)

// CTX key
type contextKey int

// GO routine response
type Response struct {
	err error
}

// Process type function signature
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

// Ternary operator
func If[T any](cond bool, trueReturn, falseReturn T) T {
	if cond {
		return trueReturn
	}
	return falseReturn
}

/*
Entry Point
- Merkletree service configuration setup and start of request.
*/
func DeriveRoot(hashes [][]byte, algorithmRequested string, processType int, initEncoding bool) ([]byte, error) {
	// Check if requested algorithm is available.
	algorithmRequested = strings.ToUpper(algorithmRequested)
	if _, ok := AlgorithmRegistry[algorithmRequested]; !ok {
		return []byte{}, &UnknownAlgorithmErr{algorithmRequested}
	}

	// Check if we got something to work with.
	if len(hashes) == 0 {
		return []byte{}, &EmptyListErr{}
	}

	// Set/check process type request.
	if processType < PassThrough || processType > BinaryTree {
		return []byte{}, &InvalidProcessTypeErr{}
	}

	// Initialize merkle pertinents.
	ms := &MerkleServer{
		Leaves:           hashes,
		HashTypeID:       algorithmRequested,
		hashGenerator:    AlgorithmRegistry[algorithmRequested],
		ProcessType:      processType,
		InitWithEncoding: initEncoding,
	}

	// Registe process type functions
	ms.ProcessTypeRegistry = map[int]processTypeFunction{
		0: ms.processPassThroughRequest,
		1: ms.processDuplicateAndAppendRequest,
		2: ms.processBinaryTreeRequest,
	}

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
		// return []byte{}, fmt.Errorf("timed out: %+w", ctx.Err())
		return []byte{}, &ProcessTimedOutErr{ctx.Err()}
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
