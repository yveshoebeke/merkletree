package merkletree

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"sort"
)

var AlgorithmRegistry map[string]CryptoFunc

type CryptoFunc func([]byte) []byte

func NOP(hash []byte) []byte {
	return hash[:]
}

func MD5(hash []byte) []byte {
	sumResult := md5.Sum(hash)
	return sumResult[:]
}

func SHA1(hash []byte) []byte {
	sumResult := sha1.Sum(hash)
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
		"NOP":		NOP,
		"MD5":          MD5,
		"SHA1":         SHA1,
		"SHA256SUM256": SHA256SUM256,
		"SHA512SUM256": SHA512SUM256,
		"SHA512SUM512": SHA512SUM512,
	}
}

// Returns json string with all available hash functions.
func AvailableAlgorithms() (string, error) {
	type availableJson struct {
		Algorithms []string `json:"algorithms"`
	}
	jsonResult := &availableJson{}

	for algorithmName := range AlgorithmRegistry {
		jsonResult.Algorithms = append(jsonResult.Algorithms, algorithmName)
	}

	sort.Strings(jsonResult.Algorithms)
	jsonEncodedAlgorithms, err := json.Marshal(jsonResult)
	if err != nil {
		return "", err
	}

	return string(jsonEncodedAlgorithms), nil
}
