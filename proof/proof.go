package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"slices"
	"strings"
)

var (
	rawData     [][]byte
	leaves      [][]byte
	encodedData []byte
)

func removeNillBytes(leaves [][]byte) [][]byte {
	return slices.DeleteFunc(leaves, func(leaf []byte) bool {
		return len(leaf) == 0
	})
}

func encodeData(data []byte) []byte {
	eData := sha256.Sum256([]byte(data))
	return eData[:]
}

func printBranch(msg string, leaves [][]byte) {
	fmt.Printf("\n%s\n", msg)
	for i, l := range leaves {
		fmt.Printf("%d -> [%s]\n", i, hex.EncodeToString(l[:]))
	}

}

func main() {
	rawDataString := "I want proof right now"
	rawDataStringSlice := strings.Split(rawDataString, " ")
	fmt.Printf("len=%d cap=%d %v\n", len(rawData), cap(rawData), rawData)

	fmt.Printf("raw data:\n")
	for _, data := range rawDataStringSlice {
		encodedData = encodeData([]byte(data))
		fmt.Printf("[%s]\t[%s]\n", data, hex.EncodeToString(encodedData))
		rawData = append(rawData, encodedData)
	}

	fmt.Println("\n------------------------------------------------------------------------")
	fmt.Println("<if odd number of elements, pass last to next branch>")
	//passthrough
	leaves = rawData

	for len(leaves) > 1 {
		for index := 0; index < len(leaves); index += 2 {
			if index+1 >= len(leaves) {
				break
			}
			leaves[index] = append(leaves[index][:], leaves[index+1][:]...)
			printBranch("appending", leaves)
			leaves[index+1] = []byte{}
			printBranch("nilling", leaves)
			leaves[index] = encodeData(leaves[index])
			printBranch("encoding", leaves)

		}
		leaves = removeNillBytes(leaves)
		printBranch("removing nills", leaves)
	}

	fmt.Printf("\npassthrough root:\n%s\n", hex.EncodeToString(leaves[0]))

	fmt.Println("\n------------------------------------------------------------------------")
	fmt.Println("<If odd number of elements, duplicate last and append>")
	//dupeandappend
	rawData = [][]byte{}
	fmt.Printf("raw data:\n")
	for _, data := range rawDataStringSlice {
		encodedData = encodeData([]byte(data))
		fmt.Printf("[%s]\t[%s]\n", data, hex.EncodeToString(encodedData))
		rawData = append(rawData, encodedData)
	}

	leaves = rawData
	started := false
	for {
		// One remaining: Exit the loop. Merkle tree root determined.
		// Note: if initial data set is only 1 element, will continue
		//	so as to adhere to this Merkle tree discipline.
		if started && len(leaves) == 1 {
			break
		}
		started = true
		//  Adjust for odd number of leaves by duplicating last leave and appending it.
		//	- ie:
		//		[1] [2] [3] [4] [5] => [1] [2] [3] [4] [5] [5]
		printBranch("branch", leaves)
		if len(leaves)%2 == 1 {
			leaves = append(leaves, leaves[len(leaves)-1])
		}
		printBranch("branch (if odd length added last one)", leaves)

		// Create combined (concatenated) hash of left and right (in couple),
		//	transform with requested algorithm (hash), and
		// 	store it in left and zero out right.
		//	- ie:
		// 		[1] [2] [3] [4] => [12] [0] [34] [0]
		for index := 0; index < len(leaves); index += 2 {
			leaves[index] = append(leaves[index][:], leaves[index+1][:]...)
			printBranch("appending", leaves)
			leaves[index+1] = []byte{}
			printBranch("nilling", leaves)
			leaves[index] = encodeData(leaves[index])
			printBranch("encoding", leaves)

		}

		// Remove 'nill' bytes.
		leaves = removeNillBytes(leaves)
		printBranch("removing nills", leaves)

	}

	fmt.Printf("\ndupeandappend root:\n%s\n", hex.EncodeToString(leaves[0]))

	fmt.Println("\n------------------------------------------------------------------------")
	fmt.Println("<Binary Tree>")
	//binarytree
	rawData = [][]byte{}
	fmt.Printf("raw data:\n")
	for _, data := range rawDataStringSlice {
		encodedData = encodeData([]byte(data))
		fmt.Printf("[%s]\t[%s]\n", data, hex.EncodeToString(encodedData))
		rawData = append(rawData, encodedData)
	}

	leaves = rawData
	var index int
	startIndex := int(math.Pow(2, math.Ceil(math.Log2(float64(len(leaves)))))) - len(leaves)
	fmt.Println("starting index =", startIndex)

	fmt.Println("[-------\nExplaines:")
	x := len(leaves)
	fmt.Println("number of data elements:", x)
	y := int(math.Ceil(math.Log2(float64(len(leaves)))))
	z := math.Pow(2, float64(y))
	fmt.Printf("lowest int exponent value to raise 2 by, resulting in a value >= to number of elements: %d\n", int(z))
	fmt.Printf("log2 (binary logarithm) of %d is %0.5f... round it up => %d\n", x, math.Log2(float64(len(leaves))), y)
	fmt.Printf("2 raised by %d = %d, which is >= %d\n", y, int(z), x)
	fmt.Printf("that value (%d) minus number of elements (%d) is the starting index: %d\n", int(z), x, int(z-float64(x)))
	fmt.Println("-------]")

	for index = startIndex; index < len(leaves); index += 2 {
		printBranch("branch", leaves)
		leaves[index] = append(leaves[index][:], leaves[index+1][:]...)
		printBranch("appending", leaves)
		leaves[index] = encodeData(leaves[index])
		printBranch("encoding", leaves)
		leaves[index+1] = []byte{}
		printBranch("nilling", leaves)
	}

	// Removenill bytes
	leaves = removeNillBytes(leaves)
	printBranch("removing nills", leaves)

	fmt.Println("\n-> now we have a binary tree")
	for len(leaves) > 1 {
		for index = 0; index < len(leaves); index += 2 {
			printBranch("branch", leaves)
			leaves[index] = append(leaves[index][:], leaves[index+1][:]...)
			printBranch("appending", leaves)
			leaves[index] = encodeData(leaves[index])
			printBranch("encoding", leaves)
			leaves[index+1] = []byte{}
			printBranch("nilling", leaves)
		}

		// Removenill bytes
		leaves = removeNillBytes(leaves)
		printBranch("removing nills", leaves)
	}

	fmt.Printf("\nbinarytree root:\n%s\n", hex.EncodeToString(leaves[0]))
}
