package nandcpu

import (
	"bufio"
	"regexp"
	"os"
	)

// Generic token container
type Token struct {
	Match string
	Type uint
}

// internal errors
var (
	ENOMATCHER = os.NewError("No valid matcher found")
	EREAD = os.NewError("Error while reading")
	)

type Matcher regexp.Regexp

// Identifier for different token types
const (
	IDENTIFIER = iota
	NUMBER
	PLUS
	MINUS
	EQUALS
	COLON
	)


// globals
var (
	tokenMatcher = []*Matcher {
		(*Matcher)(regexp.MustCompile("^[\\-.a-zA-Z0-9]+$")),
	}
	)

func (m *Matcher) Match(data []byte) bool {
	r := (*regexp.Regexp)(m)
	return r.Match(data)
}


type Lexer struct {
	input *bufio.Reader
}

func NewLexer(input *bufio.Reader) *Lexer {
	return &Lexer{input: input,
	}
}

// Returns the next token from the input
func (l *Lexer) GetNextToken() (*Token, os.Error) {
	chunk := make([]byte, 0)
	ttype, e := l.findMatcher(&chunk)
	if e != nil {
		return nil, e
	}
	e = l.maximizeMatch(ttype, &chunk)
	if e != nil {
		return nil, e
	}
	t := new(Token)
	t.Match = string(chunk)
	t.Type = ttype
	l.Reset()
	return t, nil
}

// 
func (l *Lexer) Reset() {
}

// Adds characters to the match until the matcher doesn't accept
// anymore, unreads the last byte and returns the remaining chunk
func (l *Lexer) maximizeMatch(ttype uint, chunk *[]byte) os.Error {
	println("Start maximizing")
	matcher := tokenMatcher[ttype]
	for matcher.Match(*chunk) {
		println("Chunk: ", string(*chunk))
		n, e := l.input.Read(appendEmptyCell(chunk))
		if n != 1 || e != nil {
			return EREAD
		}
	}
	l.input.UnreadByte()
	*chunk = (*chunk)[0:len(*chunk)-1]
	return nil
}

// Finds the first matcher which accepts the input
// If multiple matcher match, the one with the lower array 
// index will be used
func (l *Lexer) findMatcher(chunk *[]byte) (uint, os.Error) {
	for {
		n, e := l.input.Read(appendEmptyCell(chunk))
		if n != 1 || e != nil {
			return 0, e
		}
		ttype, e := findTokenType(*chunk)
		if e == nil {
			println("Found a matcher: ", ttype)
			return ttype, nil
		}
	}
	panic("I escaped an infinite loop... FREEDOM!")
}

// Test all tokenMatcher and returns the type index of
// the token
func findTokenType(chunk []byte) (uint, os.Error) {
	for i, matcher := range tokenMatcher {
		if(matcher.Match(chunk)) {
			return uint(i), nil
		}
	}
	return 0, ENOMATCHER
}

// Takes a slice, grows it by one and returns
// a slice with the new field only
func appendEmptyCell(slice *[]byte) []byte {
	*slice = append(*slice, 0)
	return (*slice)[len(*slice)-1:len(*slice)]
}