package stack

import (
	"context"
	"errors"
)

// ErrIsEmpty is returned when stack is empty.
var ErrIsEmpty = errors.New("stack is empty")

// Stack represents stack.
type Stack interface {
	Push(ctx context.Context, val any) error
	Pop(ctx context.Context) (any, error)
	Peek(ctx context.Context) (any, error)
}
