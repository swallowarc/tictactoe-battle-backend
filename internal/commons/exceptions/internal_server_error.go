package exceptions

import (
	"golang.org/x/xerrors"
)

// InternalError is Other errors.
//   サーバ側起因で発生したシステムエラー.
//   リクエストする側から原因が分からないため、リクエストの内容に起因するものには使用しないこと.
type InternalServerError struct {
	error
}

func IsInternalServerError(err error) bool {
	return xerrors.As(err, &InternalServerError{})
}

func NewInternalServerError(text string) InternalServerError {
	return InternalServerError{error: xerrors.New(text)}
}

func NewInternalServerErrorWrapErr(text string, err error) InternalServerError {
	e := xerrors.Errorf(text+": %w", err)
	return InternalServerError{error: e}
}
