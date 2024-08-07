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
type MerkleService struct {
	Leaves              [][]byte                    `json:"-"`
	HashTypeID          string                      `json:"hashtype"`
	hashGenerator       CryptoFunc                  `json:"-"`
	ProcessType         int                         `json:"processtype"`
	ProcessTypeRegistry map[int]processTypeFunction `json:"-"`
	ProofRequest        bool                        `json:"proofrequest"`
	ProcessResult       []byte                      `json:"root"`
	ProofResult         []byte                      `json:"proofresult"`
}

/*
Entry Point
- Merkletree service configuration setup and start of request.
*/
func DeriveRoot(hashes [][]byte, algorithmRequested string, processType int) ([]byte, error) {
	// Validate arguments
	if err := validateArgs(hashes, algorithmRequested, processType); err != nil {
		return []byte{}, err
	}

	// Initialize merkle pertinents.
	ms := &MerkleService{
		Leaves:        hashes,
		HashTypeID:    algorithmRequested,
		hashGenerator: AlgorithmRegistry[algorithmRequested],
		ProcessType:   processType,
	}

	// Registe process type functions
	ms.ProcessTypeRegistry = map[int]processTypeFunction{
		0: ms.processPassThroughRequest,
		1: ms.processDuplicateAndAppendRequest,
		2: ms.processBinaryTreeRequest,
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
		return []byte{}, &ProcessTimedOutErr{ctx.Err()}
	case resp := <-resch:
		if resp.err != nil {
			return []byte{}, resp.err
		}
	}

	return ms.ProcessResult, nil
}

// Arguments validation
func validateArgs(data [][]byte, algoReq string, pType int) error {
	var (
		validationErrs []string
		sb             strings.Builder
	)

	// check if we got something to work with.
	if len(data) == 0 {
		validationErrs = append(validationErrs, "empty data")
	}
	// Validate existing algorithm request
	if _, ok := AlgorithmRegistry[strings.ToUpper(algoReq)]; !ok {
		validationErrs = append(validationErrs, "unknown algorithm")
	}
	// is process type within range
	if pType < PassThrough || pType > BinaryTree {
		validationErrs = append(validationErrs, "invalid process type")
	}
	// nothing detected: retrun nil
	if len(validationErrs) == 0 {
		return nil
	}

	// construct error message and return it
	for _, valErr := range validationErrs {
		sb.WriteString(fmt.Sprintf("%s - ", valErr))
	}

	return &ArgumentErr{sb.String()}
}

// Ternary operator
func If[T any](cond bool, trueReturn, falseReturn T) T {
	if cond {
		return trueReturn
	}
	return falseReturn
}

// Remove elements with nill byte content and collapse slice.
//   - ie:
//     [12] [nil] [34] [nil] [56] [78] => [12] [34] 56] [78]
func (ms *MerkleService) removeNillBytes(processType int, startValue ...int) {
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
