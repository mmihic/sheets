package formula

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
)

// ParseErrorf returns a formula error.
func ParseErrorf(pos lexer.Position, msg string, args ...any) error {
	return WrapParseError(pos, fmt.Errorf(msg, args...))
}

// WrapParseError wraps an error in a parse error.
func WrapParseError(pos lexer.Position, err error) error {
	return fmt.Errorf("error at %s: %w", pos, err)
}
