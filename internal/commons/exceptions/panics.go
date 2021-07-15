package exceptions

import (
	"fmt"

	"golang.org/x/xerrors"
)

// Panic is
//   Loggerが使用できない箇所でPanicを発生させる場合にのみ使用する。
//   Loggerが使用できる場合はLogger.Fatalを使うようにする。
func Panic(message string) {
	panic(xerrors.Errorf(message))
}

// PanicWithError is
//   Loggerが使用できない箇所でPanicを発生させる場合にのみ使用する。
//   Loggerが使用できる場合はLogger.Fatalを使うようにする。
func PanicWithError(message string, v interface{}) {
	panic(xerrors.Errorf(fmt.Sprintf(message+": %+v\n"), v))
}
