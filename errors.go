package merkletree

import (
	"fmt"
)

// Custom error definitions.

// - argument validation error(s)
type ArgumentErr struct {
	invalidArguments string
}

func (argerr *ArgumentErr) Error() string {
	return fmt.Sprintf("argument error(s) - %s", argerr.invalidArguments)
}

// - process type value does not match context
type InvalidContextProcessTypeErr struct {
	contextProcess string
}

func (ctxNomatch *InvalidContextProcessTypeErr) Error() string {
	return fmt.Sprintf("process type does not match context: %s", ctxNomatch.contextProcess)
}

// - Process timed out
type ProcessTimedOutErr struct {
	ctxError error
}

func (ctxTimeout *ProcessTimedOutErr) Error() string {
	return fmt.Sprintf("timed out: %+v", ctxTimeout.ctxError)
}
