package i18n

import "fmt"

// Errorf localizes the format string and then behaves exactly like fmt.Errorf
// (including %w error wrapping). Because the format parameter is forwarded to
// fmt.Errorf, `go vet` treats Errorf as a printf-style wrapper and validates
// the format/verb pairing at each call site — where the format is a constant
// English string — instead of flagging the non-constant format here.
func Errorf(format string, a ...any) error {
	format = T(format)
	return fmt.Errorf(format, a...)
}
