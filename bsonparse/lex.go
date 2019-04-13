package bsonparse

import (
	"bytes"
	"errors"
	"strconv"
	"unicode"
)

//go:generate goyacc -l -o parser.go parser.y

// Parse parses the input and returs the result.
func Parse(input []byte) ([]interface{}, error) {
	l := newLex(input)
	_ = yyParse(l)
	return l.result, l.err
}

type lex struct {
	input  []byte
	pos    int
	result []interface{}
	err    error
}

func newLex(input []byte) *lex {
	return &lex{
		input: input,
	}
}

// Lex satisfies yyLexer.
func (l *lex) Lex(lval *yySymType) int {
	return l.scanNormal(lval)
}

func isLetter(b byte) bool {
	return unicode.IsLetter(rune(b)) ||
		unicode.IsDigit(rune(b)) ||
		b == '_' ||
		b == '+' || b == '-'
}

func isSpecialChar(b byte) bool {
	return b == '[' || b == ']' ||
		b == '{' || b == '}' ||
		b == '(' || b == ')' ||
		b == ':' || b == 0 || b == ',' || b == '"'
}

func isSpecialCharNoQuote(b byte) bool {
	return b == '[' || b == ']' ||
		b == '{' || b == '}' ||
		b == '(' || b == ')' ||
		b == ':' || b == 0 || b == ','
}

func (l *lex) scanNormal(lval *yySymType) int {
	for b := l.next(); b != 0; b = l.next() {
		switch {
		case unicode.IsSpace(rune(b)):
			continue
		case isSpecialChar(b):
			return int(b)
		default:
			l.backup()
			result := l.scanAll(lval)
			return result
		}
	}
	return 0
}

var escape = map[byte]byte{
	'"':  '"',
	'\\': '\\',
	'/':  '/',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
}

var literal = map[string]interface{}{
	"true":  true,
	"false": false,
	"null":  nil,
}

func (l *lex) scanAll(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	l.backup()
	first := l.next()
	for b := l.next(); b != 0; b = l.next() {
		if b == '\\' {
			// TODO: handle \uxxxx construct.
			b2 := escape[l.next()]
			if b2 == 0 {
				return LexError
			}
			buf.WriteByte(b2)
		} else if b == 0 ||
			(unicode.IsSpace(rune(first)) && isSpecialCharNoQuote(b)) ||
			(isSpecialCharNoQuote(first) && isSpecialCharNoQuote(b)) ||
			(first == '"' && b == '"') {
			l.backup()
			currentStr := buf.String()
			val, ok := literal[currentStr]
			if !ok {
				if currentStr == "undefined" {
					return Undefined
				} else if currentStr == "MinKey" {
					return MinKey
				} else if currentStr == "MaxKey" {
					return MaxKey
				} else if currentStr == "ObjectId" && b == '(' {
					return ObjectID
				} else if currentStr == "ISODate" && b == '(' {
					return ISODate
				} else if currentStr == "NumberLong" && b == '(' {
					return NumberLong
				} else if currentStr == "NumberDecimal" && b == '(' {
					return NumberDecimal
				} else if currentStr == "DBRef" && b == '(' {
					return DBRef
				} else {
					if b == '"' {
						lval.val = currentStr
						return String
					}
					val, err := strconv.ParseFloat(currentStr, 64)
					if err != nil {
						lval.val = currentStr
						return String
					}
					lval.val = val
					return Number
				}
			}
			lval.val = val
			return Literal
		} else if unicode.IsSpace(rune(b)) {
			continue
		} else {
			buf.WriteByte(b)
		}
	}
	return LexError
}

func (l *lex) backup() {
	if l.pos == -1 {
		return
	}
	l.pos--
}

func (l *lex) next() byte {
	if l.pos >= len(l.input) || l.pos == -1 {
		l.pos = -1
		return 0
	}
	l.pos++
	return l.input[l.pos-1]
}

// Error satisfies yyLexer.
func (l *lex) Error(s string) {
	l.err = errors.New(s)
}
