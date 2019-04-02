package bsonparse

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"unicode"
)

//go:generate goyacc -l -o parser.go parser.y

// Parse parses the input and returs the result.
func Parse(input []byte) (map[string]interface{}, error) {
	l := newLex(input)
	_ = yyParse(l)
	return l.result, l.err
}

type lex struct {
	input  []byte
	pos    int
	result map[string]interface{}
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

func (l *lex) scanNormal(lval *yySymType) int {
	for b := l.next(); b != 0; b = l.next() {
		switch {
		case unicode.IsSpace(rune(b)):
			continue
		case b == '"':
			return l.scanString(lval)
		case unicode.IsDigit(rune(b)) || b == '+' || b == '-':
			l.backup()
			return l.scanNum(lval)
		case unicode.IsLetter(rune(b)):
			l.backup()
			literal := l.scanLiteral(lval)
			if literal == LexError {
				l.backup()
				return l.scanObj(lval)
			}
			return literal
		default:
			return int(b)
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

func (l *lex) scanString(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for b := l.next(); b != 0; b = l.next() {
		switch b {
		case '\\':
			// TODO(sougou): handle \uxxxx construct.
			b2 := escape[l.next()]
			if b2 == 0 {
				return LexError
			}
			buf.WriteByte(b2)
		case '"':
			lval.val = buf.String()
			return String
		default:
			buf.WriteByte(b)
		}
	}
	return LexError
}

func (l *lex) scanNum(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for {
		b := l.next()
		switch {
		case unicode.IsDigit(rune(b)):
			buf.WriteByte(b)
		case strings.IndexByte(".+-eE", b) != -1:
			buf.WriteByte(b)
		default:
			l.backup()
			val, err := strconv.ParseFloat(buf.String(), 64)
			if err != nil {
				return LexError
			}
			lval.val = val
			return Number
		}
	}
}

var literal = map[string]interface{}{
	"true":  true,
	"false": false,
	"null":  nil,
}

func (l *lex) scanLiteral(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	counter := 0
	for {
		b := l.next()
		switch {
		case unicode.IsLetter(rune(b)):
			buf.WriteByte(b)
		default:
			l.backup()
			val, ok := literal[buf.String()]
			if !ok {
				if buf.String() == "undefined" {
					return Undefined
				} else if buf.String() == "MinKey" {
					return MinKey
				} else if buf.String() == "MaxKey" {
					return MaxKey
				} else {
					for i := 0; i < counter-1; i++ {
						l.backup()
					}
					return LexError
				}
			}
			lval.val = val
			return Literal
		}
		counter++
	}
}

func (l *lex) scanObj(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	objBuf := bytes.NewBuffer(nil)
	isInside := false
	for {
		b := l.next()
		if b == '"' && !isInside {
			isInside = true
		} else if b == '"' && isInside {
			// l.next()
			lval.val = buf.String()
			if objBuf.String() == "ObjectId(" && len(buf.String()) == 24 {
				return ObjectID
			} else if objBuf.String() == "ISODate(" {
				return ISODate
			} else if objBuf.String() == "NumberLong(" {
				return NumberLong
			} else if objBuf.String() == "NumberDecimal(" {
				return NumberDecimal
			} else if objBuf.String() == "DBRef(" {
				return DBRef
			} else {
				return BsonError
			}
		} else if isInside {
			buf.WriteByte(b)
		} else {
			objBuf.WriteByte(b)
		}
	}
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
