package sheets

import (
	"fmt"
	"strings"

	"github.com/mmihic/sheets/src/pkg/sheets/internal/formula"
)

const (
	enableParseTracing = false
)

// ParseFormula parses a formula.
//
// Formulas don't really follow a context-free grammar, but a pseudo-EBNF
// looks something like:
//
// Formula 			:= <Expression> { ">" | "<" | ">=" | "<=" | "<>" <Expression> }
// Expression 		:= <Term> { ("+" | "-") <Term> }*
// Term       		:= <Factor>  { ("*"|"/") <Expression>}*
// Factor     		:= <FunctionCall> | <Reference> | <Constant> | "(" <Formula> ")"
// FunctionCall		:= IDENTIFIER '(' <ArgList>? ')'
// ArgList			:= <Formula> (',' <Formula>)*
// Reference		:= <Sheet>? (CELL | CELL_RANGE | NAMED_RANGE)
// Constant			:= STRING | NUMBER | TRUE | FALSE
// STRING			= QuotedString
// TRUE				= [Tt][Rr][Uu][Ee]
// FALSE			= [Ff][Aa][Ll][Ss][Ee]
// IDENTIFIER		= [A-Aa-z_][A-Za-z0-9_]*
// CELL 			= ([A-Za-z]+)(0-9+)
// CELL_RANGE		= ([A-Za-z]+)?([0-9]+)?\s*:\s*([A-Za-z]+)?([0-9]+)?
func ParseFormula(s string) (Formula, error) {
	lex, err := formula.LexString(s)
	if err != nil {
		return nil, err
	}

	f, err := parseFormula(lex, 0)
	if err != nil {
		return nil, err
	}

	// Check final token is EOF
	if last, err := lex.Next(); err != nil {
		return nil, err
	} else if !last.EOF() {
		return nil, unexpectedTokenError(last, "EOF")
	}
	return f, nil
}

func parseFormula(lex *formula.Lexer, depth int) (Formula, error) {
	traceParse(depth, "parsing formula")

	expr, err := parseExpression(lex, depth+1)
	if err != nil {
		return nil, err
	}

	next, err := lex.Next()
	if err != nil {
		return nil, err
	}

	switch next.Type {
	case ">", "<", ">=", "<=", "<>", "=":
		nextExpr, err := parseExpression(lex, depth+1)
		if err != nil {
			return nil, err
		}

		return &Expression{
			Left:     expr,
			Right:    nextExpr,
			Operator: Operator(next.Type),
		}, nil
	default:
		lex.Push(next)
		return expr, nil
	}
}

func parseExpression(lex *formula.Lexer, depth int) (Formula, error) {
	traceParse(depth, "parsing expression")

	term, err := parseTerm(lex, depth+1)
	if err != nil {
		return nil, err
	}

	next, err := lex.Next()
	if err != nil {
		return nil, err
	}

	switch next.Type {
	case "+", "-":
		nextTerm, err := parseTerm(lex, depth+1)
		if err != nil {
			return nil, err
		}

		return &Expression{
			Operator: Operator(next.Type),
			Left:     term,
			Right:    nextTerm,
		}, nil
	default:
		lex.Push(next)
		return term, nil
	}
}

func parseTerm(lex *formula.Lexer, depth int) (Formula, error) {
	traceParse(depth, "parsing term")

	factor, err := parseFactor(lex, depth+1)
	if err != nil {
		return nil, err
	}

	next, err := lex.Next()
	if err != nil {
		return nil, err
	}

	switch next.Type {
	case "*", "/", "^":
		nextFactor, err := parseFactor(lex, depth+1)
		if err != nil {
			return nil, err
		}

		return &Expression{
			Operator: Operator(next.Type),
			Left:     factor,
			Right:    nextFactor,
		}, nil
	default:
		lex.Push(next)
		return factor, nil
	}
}

func parseFactor(lex *formula.Lexer, depth int) (Formula, error) {
	traceParse(depth, "parsing factor")

	tok, err := lex.Next()
	if err != nil {
		return nil, err
	}

	switch tok.Type {
	case formula.TokenTypeIdent:
		// Might be a function call or a reference. Which one depends on the next token - we
		// know it is a function call if the next token is a start paren
		nextTok, err := lex.Next()
		if err != nil {
			return nil, err
		}

		lex.Push(nextTok, tok)
		if nextTok.Type == "(" {
			return parseFunction(lex, depth+1)
		}

		return parseReference(lex, depth+1)

	case formula.TokenTypeString:
		// Might be a sheet reference or a string constant. Which one depends on the next token -
		// if the next token is a ! then it's a sheet reference
		nextTok, err := lex.Next()
		if err != nil {
			return nil, err
		}

		lex.Push(nextTok, tok)
		if nextTok.Type == "!" {
			return parseReference(lex, depth+1)
		}

		return parseConstant(lex, depth+1)

	case formula.TokenTypeCellRange:
		lex.Push(tok)
		return parseReference(lex, depth+1)

	case formula.TokenTypeNumber, formula.TokenTypeTrue, formula.TokenTypeFalse:
		lex.Push(tok)
		return parseConstant(lex, depth+1)

	case "(":
		f, err := parseFormula(lex, depth+1)
		if err != nil {
			return nil, err
		}
		next, err := lex.Next()
		if err != nil {
			return nil, err
		}
		if next.Type != ")" {
			return nil, unexpectedTokenError(next, ")")
		}
		return f, nil

	default:
		return nil, unexpectedTokenError(tok,
			formula.TokenTypeIdent, formula.TokenTypeCellRange, formula.TokenTypeNumber,
			formula.TokenTypeString, formula.TokenTypeTrue, formula.TokenTypeFalse)
	}
}

func parseFunction(lex *formula.Lexer, depth int) (Formula, error) {
	traceParse(depth, "parsing function")

	fnameToken, err := lex.Next()
	if err != nil {
		return nil, err
	}

	if fnameToken.Type != formula.TokenTypeIdent {
		return nil, unexpectedTokenError(fnameToken, formula.TokenTypeIdent)
	}

	fname := strings.ToUpper(fnameToken.Value)

	// Start the argument list
	if startParen, err := lex.Next(); err != nil {
		return nil, err
	} else if startParen.Type != "(" {
		return nil, unexpectedTokenError(startParen, "(")
	}

	// Quick check to see if the argument list is empty
	if maybeEndParen, err := lex.Next(); err != nil {
		return nil, err
	} else if maybeEndParen.Type == ")" {
		// Empty argument list
		return &FunctionCall{
			FunctionName: fname,
		}, nil
	} else {
		lex.Push(maybeEndParen)
	}

	// Parse argument list
	var args []Formula
argParsingLoop:
	for {
		arg, err := parseFormula(lex, depth+1)
		if err != nil {
			return nil, err
		}

		args = append(args, arg)
		next, err := lex.Next()
		if err != nil {
			return nil, err
		}

		switch next.Type {
		case ",":
			// Move to the next parameter
			continue argParsingLoop
		case ")":
			// This is the end of the argument list
			break argParsingLoop
		default:
			return nil, unexpectedTokenError(next, ",", ")")
		}
	}

	return &FunctionCall{
		FunctionName: fname,
		Args:         args,
	}, nil
}

func parseConstant(lex *formula.Lexer, depth int) (Formula, error) {
	traceParse(depth, "parsing constant")

	tok, err := lex.Next()
	if err != nil {
		return nil, err
	}

	switch tok.Type {
	case formula.TokenTypeTrue, formula.TokenTypeFalse, formula.TokenTypeString, formula.TokenTypeNumber:
		return &Constant{
			Value: StringToValue(tok.Value),
		}, nil

	default:
		return nil, unexpectedTokenError(tok,
			formula.TokenTypeString, formula.TokenTypeNumber,
			formula.TokenTypeTrue, formula.TokenTypeFalse)
	}
}

func parseReference(lex *formula.Lexer, depth int) (Formula, error) {
	traceParse(depth, "parsing reference")

	tok, err := lex.Next()
	if err != nil {
		return nil, err
	}

	switch tok.Type {
	case formula.TokenTypeCellRange:
		r, err := ParseRange(tok.Value)
		if err != nil {
			return nil, formula.ParseErrorf(tok.Position, "invalid range: %s", err)
		}
		return &CellRangeReference{
			Range: r,
		}, nil

	case formula.TokenTypeIdent, formula.TokenTypeString:
		// Could be a sheet, a cell, or a named range
		nextTok, err := lex.Next()
		if err != nil {
			return nil, err
		}

		if nextTok.Type == "!" {
			// First token is the name of the sheet, next token must be a cell position
			sheetName := tok.Value
			return parseCellOrNamedRange(sheetName, lex, depth+1)
		}

		// The current token is either a cell or a named range
		lex.Push(nextTok)
		lex.Push(tok)
		return parseCellOrNamedRange("", lex, depth+1)

	default:
		return nil, unexpectedTokenError(tok, formula.TokenTypeIdent, formula.TokenTypeString, formula.TokenTypeCellRange)
	}
}

func parseCellOrNamedRange(sheetName string, lex *formula.Lexer, depth int) (Formula, error) {
	traceParse(depth, "parsing cell/named range")

	nextTok, err := lex.Next()
	if err != nil {
		return nil, err
	}

	switch nextTok.Type {
	case formula.TokenTypeCellRange:
		r, err := ParseRange(nextTok.Value)
		if err != nil {
			return nil, formula.WrapParseError(nextTok.Position, err)
		}

		return &CellRangeReference{
			Sheet: sheetName,
			Range: r,
		}, nil
	case formula.TokenTypeIdent:
		if pos, err := ParsePos(nextTok.Value); err == nil {
			return &CellReference{
				Sheet: sheetName,
				Pos:   pos,
			}, nil
		}

		return &NamedRangeReference{
			NamedRange: nextTok.Value,
		}, nil
	default:
		return nil, unexpectedTokenError(nextTok, formula.TokenTypeIdent, formula.TokenTypeCellRange)
	}
}

func unexpectedTokenError(tok formula.Token, expected ...string) error {
	return formula.ParseErrorf(tok.Position, "expected one of [%s]: found '%s' (%s)",
		strings.Join(append([]string{}, expected...), ", "),
		tok.Value, tok.Type)
}

func traceParse(depth int, msg string, args ...any) {
	if !enableParseTracing {
		return
	}

	fmt.Printf("%s %s", strings.Repeat(" ", depth), fmt.Sprintf(msg, args...))
}
