package merkletree

//  - If:
//		General ternary operator.
//

const (
	NopProcess  = -1
	PassThrough = 0
	DupeAppend  = 1
	BinaryTree  = 2
)

// Ternary operator
func If[T any](cond bool, trueValue, falseValue T) T {
	if cond {
		return trueValue
	}
	return falseValue
}
