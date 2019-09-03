package extras

import (
	"fmt"
	"github.com/sasbury/lg"
)

// BranchingError holds a list of child errors
type BranchingError struct {
	Children []error
}

func (be BranchingError) Error() string {
	return fmt.Sprintf("branching error with %d children", len(be.Children))
}

// BranchingAppender holds multiple appenders and calls each one in order.
// If an error occurs, a BranchingError is returned with each child error
type BranchingAppender struct {
	branches []lg.LogAppender
}

// NewBranchingAppender returns a new appender with the provided branches
func NewBranchingAppender(appenders ...lg.LogAppender) *BranchingAppender {
	return &BranchingAppender{
		branches: append([]lg.LogAppender{}, appenders...),
	}
}

// Log is the branching appenders implementation of a LogAppender
func (ba *BranchingAppender) Log(entry string) error {
	var errors []error

	for _, a := range ba.branches {
		err := a(entry)
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return BranchingError{
			Children: errors,
		}
	}

	return nil
}
