package merkletree

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"strings"
	"testing"
)

const (
	ColorGreen   = "\u001b[32m"
	ColorDefault = "\u001b[00m"
)

// table driven tests
type DeriveRootTest struct {
	arg1     [][]byte
	arg2     string
	arg3     int
	expected []byte
	err      error
}

var (
	showTestResults                                                  = flag.Bool("detail", false, "show detail")
	processName                                                      = []string{"Duplicate and Append", "Pass-Through", "Binary Tree"}
	tenThousandElements0, tenThousandElements1, tenThousandElements2 [][]byte
	iWantProofRightNow0, iWantProofRightNow1, iWantProofRightNow2    [][]byte
	DeriveRootTests                                                  []DeriveRootTest
	resultSHA256SUM256_0, resultSHA256SUM256_1, resultSHA256SUM256_2 []byte
)

func makeTestHashes() {
	resultSHA256SUM256_0, _ = hex.DecodeString("7c6aed8b3a1f18ad143ef5911302666c758a41d1588141e46a08e44561bcd582")
	resultSHA256SUM256_1, _ = hex.DecodeString("e0934a80a459a7c2256d7eef4a819d4be23b12ed59147acb981fe2e65ecc97db")
	resultSHA256SUM256_2, _ = hex.DecodeString("84d54d7074b373c94fd43e8fb1d78b7fd1925aadff0f2bf90ef1c66d5462f24f")

	for _, word := range strings.Split("I want proof right now", " ") {
		iWantProofRightNow0 = append(iWantProofRightNow0, SHA256SUM256([]byte(word)))
		iWantProofRightNow1 = append(iWantProofRightNow1, SHA256SUM256([]byte(word)))
		iWantProofRightNow2 = append(iWantProofRightNow2, SHA256SUM256([]byte(word)))
	}

	for i := 0; i < 2000; i++ {
		tenThousandElements0 = append(tenThousandElements0, iWantProofRightNow0...)
		tenThousandElements1 = append(tenThousandElements1, iWantProofRightNow0...)
		tenThousandElements2 = append(tenThousandElements2, iWantProofRightNow0...)
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
			arg1:     iWantProofRightNow0,
			arg2:     "SHA256SUM256",
			arg3:     DupeAppend,
			expected: resultSHA256SUM256_0,
			err:      nil,
		},
		DeriveRootTest{
			arg1:     iWantProofRightNow1,
			arg2:     "SHA256SUM256",
			arg3:     PassThrough,
			expected: resultSHA256SUM256_1,
			err:      nil,
		},
		DeriveRootTest{
			arg1:     iWantProofRightNow2,
			arg2:     "SHA256SUM256",
			arg3:     BinaryTree,
			expected: resultSHA256SUM256_2,
			err:      nil,
		},
	)

	for _, test := range DeriveRootTests {
		output, err := DeriveRoot(test.arg1, test.arg2, test.arg3)

		// provided by the `-detail` flag (see above)
		if *showTestResults {
			fmt.Printf("%s\n- algorithm: %s\n-   process: %s\n------- got: %s\n-- expected: %s\n----- error: %v\n%s", ColorGreen, test.arg2, processName[test.arg3], hex.EncodeToString(output), hex.EncodeToString(test.expected), err, ColorDefault)
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
		DeriveRoot(tenThousandElements0, "SHA256SUM256", DupeAppend)
	}
}

func BenchmarkDeriveRoot10000LeavesSHA256SUM256PassThrough(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DeriveRoot(tenThousandElements1, "SHA256SUM256", PassThrough)
	}
}

func BenchmarkDeriveRoot10000LeavesSHA256SUM256BinaryTree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DeriveRoot(tenThousandElements2, "SHA256SUM256", BinaryTree)
	}
}
