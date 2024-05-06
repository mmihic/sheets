package formula

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type expectedToken struct {
	Type  string
	Value string
}

func TestLexer(t *testing.T) {
	for _, tt := range []struct {
		input    string
		expected []expectedToken
	}{
		{`"this is a quoted string"`, []expectedToken{
			{"String", "this is a quoted string"},
		}},
		{`"this is a \"quoted\" string"`, []expectedToken{
			{"String", `this is a "quoted" string`},
		}},
		{`'this is a \'single quoted\' string'`, []expectedToken{
			{"String", `this is a 'single quoted' string`},
		}},
		{"`this is a tick quoted string`", []expectedToken{
			{"String", "this is a tick quoted string"},
		}},
		{"ThisIsAnIdentifier", []expectedToken{
			{"Ident", "ThisIsAnIdentifier"},
		}},
		{"TRUE", []expectedToken{
			{"True", "TRUE"},
		}},
		{"false", []expectedToken{
			{"False", "false"},
		}},
		{"AA2:B14", []expectedToken{
			{"CellRange", "AA2:B14"},
		}},
		{"`SomeSheet`!AA2:B14", []expectedToken{
			{"String", "SomeSheet"},
			{"!", "!"},
			{"CellRange", "AA2:B14"},
		}},
		{"VLOOKUP ( A:Z, A2, 1, false ) ", []expectedToken{
			{"Ident", "VLOOKUP"},
			{"(", "("},
			{"CellRange", "A:Z"},
			{",", ","},
			{"Ident", "A2"},
			{",", ","},
			{"Number", "1"},
			{",", ","},
			{"False", "false"},
			{")", ")"},
		}},
		{"((A25+15) / 79) * ( (14+B64)-C23)^2 )", []expectedToken{
			{"(", "("},
			{"(", "("},
			{"Ident", "A25"},
			{"+", "+"},
			{"Number", "15"},
			{")", ")"},
			{"/", "/"},
			{"Number", "79"},
			{")", ")"},
			{"*", "*"},
			{"(", "("},
			{"(", "("},
			{"Number", "14"},
			{"+", "+"},
			{"Ident", "B64"},
			{")", ")"},
			{"-", "-"},
			{"Ident", "C23"},
			{")", ")"},
			{"^", "^"},
			{"Number", "2"},
			{")", ")"},
		}},
		{`100.3*17 + 45 >= A34`, []expectedToken{
			{Type: "Number", Value: "100.3"},
			{Type: "*", Value: "*"},
			{Type: "Number", Value: "17"},
			{Type: "+", Value: "+"},
			{Type: "Number", Value: "45"},
			{Type: ">=", Value: ">="},
			{Type: "Ident", Value: "A34"},
		}},
	} {
		t.Run(tt.input, func(t *testing.T) {
			l, err := LexString(tt.input)
			if !assert.NoError(t, err) {
				return
			}

			var tokens []expectedToken
			for {
				tok, err := l.Next()
				if !assert.NoError(t, err) {
					return
				}

				if tok.EOF() {
					break
				}

				tokens = append(tokens, expectedToken{
					Type:  tok.Type,
					Value: tok.Value,
				})
			}

			assert.Equal(t, tt.expected, tokens)
		})
	}
}
