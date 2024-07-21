package merkletree

//  - If:
//		General ternary operator.
//

const (
	NopProcess              = -1
	PassThrough             = 0
	DupeAppend              = 1
	BinaryTree              = 2
	ProcessTimeoutMilliSecs = 10
)

// CTX key values
var processTypes = [3]string{"PAS-THRU", "DUP-APND", "BIN-TREE"}

// Ternary operator
func If[T any](cond bool, trueValue, falseValue T) T {
	if cond {
		return trueValue
	}
	return falseValue
}
