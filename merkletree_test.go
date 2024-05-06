package merkletree

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"strings"
	"testing"
)

const (
	ColorGreen   = "\u001b[32m"
	ColorDefault = "\u001b[00m"
	PassThrough  = 0
	DupeAppend   = 1
	BinaryTree   = 2
)

// table driven tests
type DeriveRootTest struct {
	arg1     string
	arg2     [][]byte
	arg3     int
	arg4     bool
	expected []byte
	err      error
}

var (
	showTestResults                                                  = flag.Bool("detail", false, "show detail")
	processName                                                      = []string{"Duplicate and Append", "Pass-Through", "Binary Tree"}
	iWantProofRightNow                                               [][]byte
	tenThousandElements                                              [][]byte
	DeriveRootTests                                                  []DeriveRootTest
	resultSHA256SUM256_0, resultSHA256SUM256_1, resultSHA256SUM256_2 []byte
)

func makeTestHashes() {
	resultSHA256SUM256_0, _ = hex.DecodeString("7c6aed8b3a1f18ad143ef5911302666c758a41d1588141e46a08e44561bcd582")
	resultSHA256SUM256_1, _ = hex.DecodeString("e0934a80a459a7c2256d7eef4a819d4be23b12ed59147acb981fe2e65ecc97db")
	resultSHA256SUM256_2, _ = hex.DecodeString("84d54d7074b373c94fd43e8fb1d78b7fd1925aadff0f2bf90ef1c66d5462f24f")

	fillerValue := md5.Sum([]byte("merkletree"))
	for i := 0; i < 10000; i++ {
		tenThousandElements = append(tenThousandElements, fillerValue[:])
	}
}

func init() {
	makeTestHashes()
}

func TestDeriveRoot(t *testing.T) {
	if *showTestResults != true {
		fmt.Println("    [Note: use `-detail` flag to see more]")
	}
	DeriveRootTests = append(DeriveRootTests,
		DeriveRootTest{
			arg1:     "SHA256SUM256",
			arg2:     iWantProofRightNow,
			arg3:     DupeAppend,
			arg4:     false,
			expected: resultSHA256SUM256_0,
			err:      nil,
		},
		DeriveRootTest{
			arg1:     "SHA256SUM256",
			arg2:     iWantProofRightNow,
			arg3:     PassThrough,
			arg4:     false,
			expected: resultSHA256SUM256_1,
			err:      nil,
		},
		DeriveRootTest{
			arg1:     "SHA256SUM256",
			arg2:     iWantProofRightNow,
			arg3:     BinaryTree,
			arg4:     false,
			expected: resultSHA256SUM256_2,
			err:      nil,
		},
	)

	for _, test := range DeriveRootTests {
		iWantProofRightNow = [][]byte{}
		for _, word := range strings.Split("I want proof right now", " ") {
			iWantProofRightNow = append(iWantProofRightNow, SHA256SUM256([]byte(word)))
		}
		test.arg2 = iWantProofRightNow
		output, err := DeriveRoot(test.arg1, test.arg2, test.arg3, test.arg4)

		// provided by the `-detail` flag (see above)
		if *showTestResults {
			fmt.Printf("%s\n- algorithm: %s\n-   process: %s\n------- got: %s\n-- expected: %s\n----- error: %v\n%s", ColorGreen, test.arg1, processName[test.arg3], hex.EncodeToString(output), hex.EncodeToString(test.expected), err, ColorDefault)
		}

		if !bytes.Equal(output, test.expected) {
			t.Errorf("(out) got %q, wanted %q, error %q", hex.EncodeToString(output), hex.EncodeToString(test.expected), err)
		}

		if err != nil && err.Error() != test.err.Error() {
			t.Errorf("(err) got %q, wanted %q", err, test.err)
		}
	}

}

func BenchmarkDeriveRoot10000LeavesSHA256SUM256DupAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DeriveRoot("SHA256SUM256", tenThousandElements, DupeAppend, true)
	}
}

func BenchmarkDeriveRoot10000LeavesSHA256SUM256PassThrough(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DeriveRoot("SHA256SUM256", tenThousandElements, PassThrough, true)
	}
}
func BenchmarkDeriveRoot10000LeavesSHA256SUM256BinaryTree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DeriveRoot("SHA256SUM256", tenThousandElements, BinaryTree, true)
	}
}
