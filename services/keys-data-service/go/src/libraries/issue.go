package libraries

import (
	"errors"
	"fmt"
)

type Issues struct {
	issues []error
}

func (stack *Issues) Push(issue error) {
	if issue != nil {
		stack.issues = append(stack.issues, issue)
	}
}

func (stack *Issues) Raise(format string, args ...any) error {
	issue := fmt.Errorf(format, args...)
	stack.issues = append(stack.issues, issue)
	return issue
}

func (stack *Issues) Flush() error {
	defer stack.Clear()

	if len(stack.issues) == 0 {
		return nil
	}

	return errors.Join(stack.issues...)
}

func (stack *Issues) Empty() bool {
	return len(stack.issues) == 0
}

func (stack *Issues) Clear() {
	stack.issues = nil
}
