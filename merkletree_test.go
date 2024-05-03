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
	resultSHA256SUM256_0, _ = hex.DecodeString("e067e8bf61a7c48115f1bb351dc27d6a003a28581357a725d813c5569531d129")
	resultSHA256SUM256_1, _ = hex.DecodeString("db84b5e863ff8dda1946029537123f058db87615a5317cd4375bc0fee751a6b1")
	resultSHA256SUM256_2, _ = hex.DecodeString("dbff0bb54efc926fc56ec7de468b173a95e6f9bc81d05292906c214ef0f0ff3a")

	for _, word := range strings.Split("I want proof right now", " ") {
		iWantProofRightNow = append(iWantProofRightNow, SHA256SUM256([]byte(word)))
	}

	fillerValue := md5.Sum([]byte("merkletree"))
	for i := 0; i < 10000; i++ {
		tenThousandElements = append(tenThousandElements, fillerValue[:])
	}

}

func init() {
	makeTestHashes()
}

func TestDeriveRoot(t *testing.T) {
	DeriveRootTests = append(DeriveRootTests,
		DeriveRootTest{
			arg1:     "SHA256SUM256",
			arg2:     iWantProofRightNow,
			arg3:     0,
			arg4:     true,
			expected: resultSHA256SUM256_0,
			err:      nil,
		},
		DeriveRootTest{
			arg1:     "SHA256SUM256",
			arg2:     iWantProofRightNow,
			arg3:     1,
			arg4:     true,
			expected: resultSHA256SUM256_1,
			err:      nil,
		},
		DeriveRootTest{
			arg1:     "SHA256SUM256",
			arg2:     iWantProofRightNow,
			arg3:     2,
			arg4:     true,
			expected: resultSHA256SUM256_2,
			err:      nil,
		},
	)

	for _, test := range DeriveRootTests {
		output, err := DeriveRoot(test.arg1, &test.arg2, test.arg3, test.arg4)

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

func BenchmarkDeriveRoot10000LeavesSHA256SUM256DupAppendInited(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DeriveRoot("SHA256SUM256", &tenThousandElements, 0, true)
	}
}

func BenchmarkDeriveRoot10000LeavesSHA256SUM256PassThroughInited(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DeriveRoot("SHA256SUM256", &tenThousandElements, 1, true)
	}
}
func BenchmarkDeriveRoot10000LeavesSHA256SUM256BinaryTreeInited(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DeriveRoot("SHA256SUM256", &tenThousandElements, 2, true)
	}
}
