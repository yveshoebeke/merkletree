package merkletree

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"strings"
	"testing"
)

// table driven tests
type DeriveRootTest struct {
	arg1     string
	arg2     [][]byte
	arg3     bool
	expected []byte
	err      error
}

const (
	ColorGreen   = "\u001b[32m"
	ColorDefault = "\u001b[00m"
)

// Note: result var names definition as follows: `result{*1}_{*2}`, where:
//
//		*1: hash algorithm intent
//		*2: number of input elements supplied
//		- Ex: resultSHA256SUM256_10000 -> byte result of 10,000 elements using sha256.Sum256
//		- Results are hex-to-string encoded for display.
//
//	Note: showTestResults:
//
//		set to true by providing the -detail flag and will display detail output for each test.
//		- Ex: `go test -detail` or `go test -v -detail`
var (
	DeriveRootTests                                    []DeriveRootTest
	ms                                                 *MerkleServer
	showTestResults                                    = flag.Bool("detail", false, "show detail")
	myHashes                                           [][]byte
	myHashes1                                          [][]byte
	myHashes3                                          [][]byte
	myHashes100                                        [][]byte
	myHashes1000                                       [][]byte
	myHashes3000                                       [][]byte
	myHashes5000                                       [][]byte
	myHashes10000                                      [][]byte
	resultAvailableAlgorithms                          string
	resultMD5_1, resultSHA1_3                          []byte
	resultSHA256SUM256_3, resultSHA512SUM512_3         []byte
	resultSHA256SUM256_10000, resultSHA512SUM256_10000 []byte
	resultSHA512SUM512_10000                           []byte
)

// create some hashed data as test input:
func hashSetup() {
	// create some data to work with as hash leaves.
	words := strings.Split("I want proof I tell you now the time has come for all good men to come to the aid of their country", " ")
	for i := 0; i < 500; i++ {
		for _, word := range words {
			myHashes = append(myHashes, MD5([]byte(word)))
		}
	}

	myHashes1 = myHashes[:1]
	myHashes3 = myHashes[:3]
	myHashes100 = myHashes[:100]
	myHashes1000 = myHashes[:1000]
	myHashes3000 = myHashes[:3000]
	myHashes5000 = myHashes[:5000]
	myHashes10000 = myHashes[:10000]

	// setup some test scenario result expectations
	resultAvailableAlgorithms = "{\"algorithms\":[\"MD5\",\"SHA1\",\"SHA256SUM256\",\"SHA512SUM256\",\"SHA512SUM512\"]}"
	resultMD5_1, _ = hex.DecodeString("6c59fac724e211816d6dcefb4a67d4bb")
	resultSHA1_3, _ = hex.DecodeString("46552928172dedc8ff8974688eda78a46bbedcfa")
	resultSHA256SUM256_3, _ = hex.DecodeString("5d5ebebb45154af61279ead547adcd8d18a8704fbacf8b77af39bdb347d30128")
	resultSHA512SUM512_3, _ = hex.DecodeString("e521735f7b7376b8bdf1980bb490bab2a0caeac2c5cacb9c368976b60c48b71711ab111025aacd19aeaafdc3684b6aa95da4ed5c021b8f1651c2163c6bc89a27")
	resultSHA256SUM256_10000, _ = hex.DecodeString("f2727ae8274d197a4546901271aa8932bc1f7e7ce3445a6673ec4cce2776aec3")
	resultSHA512SUM256_10000, _ = hex.DecodeString("89e345c09454c49f1a2769ab08ab01b8083569689f7451ce5490fbf84bfae208")
	resultSHA512SUM512_10000, _ = hex.DecodeString("b4fc7193c7e207a0457cd6427f1ef9d69da30bae8af1d12bcf0ad37e18c6765ab960c45977b23bef8c0ce47c1955f914ed8d903809862edd8cd68dfef1383276")
}

func init() {
	hashSetup()
}

// simple test to see available algorithm test
func TestAvailableAlgorithms(t *testing.T) {
	expected := resultAvailableAlgorithms
	got, err := AvailableAlgorithms()
	if *showTestResults {
		fmt.Printf("%s(Available algorithms)\n------ got: %q\n- expected: %q\n---- error: %q\n%s", ColorGreen, got, expected, err, ColorDefault)
	}
	if got != expected {
		t.Errorf("got %q, wanted %q, error: %q", got, expected, err)
	}
}

// table-driven test on DeriveRoot function
func TestDeriveRoot(t *testing.T) {
	DeriveRootTests = append(DeriveRootTests,
		DeriveRootTest{
			arg1:     "SHA256SUM256",
			arg2:     myHashes[:3],
			arg3:     false,
			expected: resultSHA256SUM256_3,
			err:      nil,
		},
	)

	for _, test := range DeriveRootTests {
		output, err := ms.DeriveRoot(test.arg1, &test.arg2, test.arg3)
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

// sha256sum256 on 100 leaves
func BenchmarkDeriveRootWith100LeavesUsingSHA256SUM256(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ms.DeriveRoot("SHA256SUM256", &myHashes100, false)
	}
}

// sha256sum256 on 1000 leaves
func BenchmarkDeriveRootWith1000LeavesUsingSHA256SUM256(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ms.DeriveRoot("SHA256SUM256", &myHashes1000, false)
	}
}

// sha256sum256 on 3000 leaves
func BenchmarkDeriveRootWith3000LeavesUsingSHA256SUM256(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ms.DeriveRoot("SHA256SUM256", &myHashes3000, false)
	}
}

// sha256sum256 on 3000 leaves
func BenchmarkDeriveRootWith5000LeavesUsingSHA256SUM256(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ms.DeriveRoot("SHA256SUM256", &myHashes5000, false)
	}
}

// sha256sum256 on 10000 leaves
func BenchmarkDeriveRootWith10000LeavesUsingSHA256SUM256(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ms.DeriveRoot("SHA256SUM256", &myHashes10000, false)
	}
}
