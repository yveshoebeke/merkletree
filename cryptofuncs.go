package merkletree

//
// Functions:
//
//	- AvailableAlgorithms (cryptofuncs.go):
//		Returns the hash algoritms available in this module.
//

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"

	"golang.org/x/crypto/sha3"
)

type CryptoFunc func([]byte) []byte

var AlgorithmRegistry map[string]CryptoFunc

func MD5(hash []byte) []byte {
	sumResult := md5.Sum(hash)
	return sumResult[:]
}

func SHA1(hash []byte) []byte {
	sumResult := sha1.Sum(hash)
	return sumResult[:]
}

func SHA3SUM256(hash []byte) []byte {
	sumResult := sha3.Sum256(hash)
	return sumResult[:]
}

func SHA256SUM256(hash []byte) []byte {
	sumResult := sha256.Sum256(hash)
	return sumResult[:]
}

func SHA512SUM256(hash []byte) []byte {
	sumResult := sha512.Sum512_256(hash)
	return sumResult[:]
}

func SHA512SUM512(hash []byte) []byte {
	sumResult := sha512.Sum512(hash)
	return sumResult[:]
}

// Create function registry
func init() {
	AlgorithmRegistry = map[string]CryptoFunc{
		"MD5":          MD5,
		"SHA1":         SHA1,
		"SHA3SUM256":   SHA3SUM256,
		"SHA256SUM256": SHA256SUM256,
		"SHA512SUM256": SHA512SUM256,
		"SHA512SUM512": SHA512SUM512,
	}
}
