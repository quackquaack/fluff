package fluff

import (
	"fmt"
	"os"
	"sync"
)

type lexer struct {
	File    *os.File
	Lexemes chan Lexeme

	Buffer [1]byte

	Line   uint64
	Column uint64
}

func (l *lexer) current() byte {
	return l.Buffer[0]
}

type LexemeType uint16

const (
	Int = iota
	Float

	Ident

	Return
)

type Lexeme struct {
	Type  LexemeType
	Value string

	Line   uint64
	Column uint64
}

func (l *lexer) next() byte {
	l.File.Read(l.Buffer[:])
	return l.current()
}

func (l *lexer) emit(ltype LexemeType, value string) {
	l.Lexemes <- Lexeme{
		Type:  ltype,
		Value: value,

		Line:   l.Line,
		Column: l.Column,
	}
}

func (l *lexer) err(msg string) string {
	return fmt.Sprint("Lexer errored at [", l.Line, ":", l.Column, "]: ", msg)
}

func skip(l *lexer) (err string) {
	if l.current() == ' ' || l.current() == '\n' || l.current() == '\t' {
		l.next()
		return ""
	}

	return ""
}

func Lex(file *os.File, lexemes chan Lexeme, wg *sync.WaitGroup) (err string) {
	defer close(lexemes)
	defer wg.Done()

	l := &lexer{
		File:    file,
		Lexemes: lexemes,

		Buffer: [1]byte{0},
	}

	for l.current() != 0 {
		err := skip(l)
		if len(err) > 0 {
			return err
		}
	}

	return ""
}
