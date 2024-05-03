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
	arg1     string
	arg2     [][]byte
	arg3     int
	arg4     bool
	expected []byte
	err      error
}

var (
	showTestResults                                                  = flag.Bool("detail", false, "show detail")
	iWantProofRightNow                                               [][]byte
	DeriveRootTests                                                  []DeriveRootTest
	resultSHA256SUM256_0, resultSHA256SUM256_1, resultSHA256SUM256_2 []byte
)

func makeTestHashes() {
	resultSHA256SUM256_0, _ = hex.DecodeString("e067e8bf61a7c48115f1bb351dc27d6a003a28581357a725d813c5569531d129")
	resultSHA256SUM256_1, _ = hex.DecodeString("64a6b27e8287e768d3f79c079517047b0fbd70423be5e9bc80fea14a5ca058c9")
	resultSHA256SUM256_2, _ = hex.DecodeString("3ccb4ab6a6d27d8c65d723345dd4d90e7d53963fb56b35fee464a999536f13e5")
}

func init() {
	makeTestHashes()
}

func TestDeriveRoot(t *testing.T) {
	DeriveRootTests = append(DeriveRootTests,
		DeriveRootTest{
			arg1:     "SHA256SUM256",
			arg2:     [][]byte{},
			arg3:     0,
			arg4:     true,
			expected: resultSHA256SUM256_0,
			err:      nil,
		},
		DeriveRootTest{
			arg1:     "SHA256SUM256",
			arg2:     [][]byte{},
			arg3:     1,
			arg4:     true,
			expected: resultSHA256SUM256_1,
			err:      nil,
		},
		DeriveRootTest{
			arg1:     "SHA256SUM256",
			arg2:     [][]byte{},
			arg3:     2,
			arg4:     true,
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
			fmt.Printf("%s\n- algorithm: %s\n------- got: %s\n-- expected: %s\n----- error: %v\n%s", ColorGreen, test.arg1, hex.EncodeToString(output), hex.EncodeToString(test.expected), err, ColorDefault)
		}

		if !bytes.Equal(output, test.expected) {
			t.Errorf("(out) got %q, wanted %q, error %q", hex.EncodeToString(output), hex.EncodeToString(test.expected), err)
		}

		if err != nil && err.Error() != test.err.Error() {
			t.Errorf("(err) got %q, wanted %q", err, test.err)
		}
	}
}
