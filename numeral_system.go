package lexorank

import (
	"fmt"
)

const (
	map36 = "0123456789abcdefghijklmnopqrstuvwxyz"
	map64 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ^_abcdefghijklmnopqrstuvwxyz"
)

type LexoNumeralSystem interface {
	GetBase() byte
	GetPositiveChar() byte
	GetNegativeChar() byte
	GetRadixPointChar() byte
	Digit(ch byte) (byte, error)
	Char(digit byte) byte
}

var (
	_ LexoNumeralSystem = (*LexoNumeralSystem10)(nil)
	_ LexoNumeralSystem = (*LexoNumeralSystem36)(nil)
	_ LexoNumeralSystem = (*LexoNumeralSystem64)(nil)
)

type LexoNumeralSystem10 struct {
}

func NewLexoNumeralSystem10() *LexoNumeralSystem10 {
	return &LexoNumeralSystem10{}
}

func (n *LexoNumeralSystem10) GetBase() byte {
	return 10
}

func (n *LexoNumeralSystem10) GetPositiveChar() byte {
	return '+'
}

func (n *LexoNumeralSystem10) GetNegativeChar() byte {
	return '-'
}

func (n *LexoNumeralSystem10) GetRadixPointChar() byte {
	return '.'
}

func (n *LexoNumeralSystem10) Digit(ch byte) (byte, error) {
	if ch >= '0' && ch <= '9' {
		return ch - '0', nil
	}

	return 0, fmt.Errorf("not valid digit: %q", ch)
}

func (n *LexoNumeralSystem10) Char(digit byte) byte {
	return digit + '0'
}

type LexoNumeralSystem36 struct {
}

func NewLexoNumeralSystem36() LexoNumeralSystem {
	return &LexoNumeralSystem36{}
}

func (n *LexoNumeralSystem36) GetBase() byte {
	return 36
}

func (n *LexoNumeralSystem36) GetPositiveChar() byte {
	return '+'
}

func (n *LexoNumeralSystem36) GetNegativeChar() byte {
	return '-'
}

func (n *LexoNumeralSystem36) GetRadixPointChar() byte {
	return ':'
}

func (n *LexoNumeralSystem36) Digit(ch byte) (byte, error) {
	switch {
	case ch >= '0' && ch <= '9':
		return ch - '0', nil

	case ch >= 'a' && ch <= 'z':
		return ch - 'a' + 10, nil
	}

	return 0, fmt.Errorf("not valid digit: %q", ch)
}

func (n *LexoNumeralSystem36) Char(digit byte) byte {
	return map36[digit]
}

type LexoNumeralSystem64 struct {
}

func NewLexoNumeralSystem64() *LexoNumeralSystem64 {
	return &LexoNumeralSystem64{}
}

func (n *LexoNumeralSystem64) GetBase() byte {
	return 64
}

func (n *LexoNumeralSystem64) GetPositiveChar() byte {
	return '+'
}

func (n *LexoNumeralSystem64) GetNegativeChar() byte {
	return '-'
}

func (n *LexoNumeralSystem64) GetRadixPointChar() byte {
	return ':'
}

func (n *LexoNumeralSystem64) Digit(ch byte) (byte, error) {
	switch {
	case ch >= '0' && ch <= '9':
		return ch - '0', nil

	case ch >= 'A' && ch <= 'Z':
		return ch - 'A' + 10, nil

	case ch == '^':
		return 36, nil

	case ch == '_':
		return 37, nil

	case ch >= 'a' && ch <= 'z':
		return ch - 'a' + 38, nil
	}

	return 0, fmt.Errorf("not valid digit: %q", ch)
}

func (n *LexoNumeralSystem64) Char(digit byte) byte {
	return map64[digit]
}
