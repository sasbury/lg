package extras

import (
	"fmt"
)

// BadAppender always returns an error with the entry as the content
func BadAppender(entry string) error {
	return fmt.Errorf(entry)
}
