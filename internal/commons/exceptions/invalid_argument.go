package exceptions

import (
	"golang.org/x/xerrors"
)

type (
	InvalidArgumentError struct {
		error
	}
)

func IsInvalidArgumentError(err error) bool {
	return xerrors.As(err, &InvalidArgumentError{})
}

func NewInvalidArgumentError(text string) InvalidArgumentError {
	return InvalidArgumentError{error: xerrors.New(text)}
}
