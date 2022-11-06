package lexorank

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexoIntegerParse(t *testing.T) {
	tests := []struct {
		name    string
		strFull string
		system  LexoNumeralSystem
		want    *LexoInteger
		wantErr bool
	}{
		{
			name:    "empty",
			strFull: "0",
			system:  NewLexoNumeralSystem36(),
			want: &LexoInteger{
				sys:  NewLexoNumeralSystem36(),
				sign: 0,
				mag:  []byte{0},
			},
			wantErr: false,
		},
		{
			name:    "zero",
			strFull: "0",
			system:  NewLexoNumeralSystem36(),
			want: &LexoInteger{
				sys:  NewLexoNumeralSystem36(),
				sign: 0,
				mag:  []byte{0},
			},
			wantErr: false,
		},
		{
			name:    "one",
			strFull: "1",
			system:  NewLexoNumeralSystem36(),
			want: &LexoInteger{
				sys:  NewLexoNumeralSystem36(),
				sign: 1,
				mag:  []byte{1},
			},
			wantErr: false,
		},
		{
			name:    "10000",
			strFull: "10000",
			system:  NewLexoNumeralSystem36(),
			want: &LexoInteger{
				sys:  NewLexoNumeralSystem36(),
				sign: 1,
				mag:  []byte{0, 0, 0, 0, 1},
			},
			wantErr: false,
		},
		{
			name:    "-10000",
			strFull: "-10000",
			system:  NewLexoNumeralSystem36(),
			want: &LexoInteger{
				sys:  NewLexoNumeralSystem36(),
				sign: -1,
				mag:  []byte{0, 0, 0, 0, 1},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LexoIntegerParse(tt.strFull, tt.system)
			if (err != nil) != tt.wantErr {
				t.Errorf(fmt.Sprintf("LexoIntegerParse(%v, %v) error = %v, wantErr %v", tt.strFull, tt.system, err, tt.wantErr))
				return
			}
			assert.Equalf(t, tt.want, got, "LexoIntegerParse(%v, %v)", tt.strFull, tt.system)
		})
	}
}

func TestLexoInteger_String(t *testing.T) {
	tests := []struct {
		name string
		i    *LexoInteger
		want string
	}{
		{
			name: "zero",
			i: &LexoInteger{
				sys:  NewLexoNumeralSystem36(),
				sign: 0,
				mag:  []byte{0},
			},
			want: "0",
		},
		{
			name: "one",
			i: &LexoInteger{
				sys:  NewLexoNumeralSystem36(),
				sign: 1,
				mag:  []byte{1},
			},
			want: "1",
		},
		{
			name: "10000",
			i: &LexoInteger{
				sys:  NewLexoNumeralSystem36(),
				sign: 1,
				mag:  []byte{0, 0, 0, 0, 1},
			},
			want: "10000",
		},
		{
			name: "-10000",
			i: &LexoInteger{
				sys:  NewLexoNumeralSystem36(),
				sign: -1,
				mag:  []byte{0, 0, 0, 0, 1},
			},
			want: "-10000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.i.String(), "String()")
		})
	}
}

func TestLexoInteger_ShiftLeft(t *testing.T) {
	tests := []struct {
		name  string
		digit string
		times int
		want  string
	}{
		{
			name:  "zero times",
			digit: "1234",
			times: 0,
			want:  "1234",
		},
		{
			name:  "one times",
			digit: "1234",
			times: 1,
			want:  "12340",
		},
		{
			name:  "more times",
			digit: "1234",
			times: 4,
			want:  "12340000",
		},
		{
			name:  "negative times",
			digit: "1234",
			times: -2,
			want:  "12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			digit, _ := LexoIntegerParse(tt.digit, NewLexoNumeralSystem36())
			assert.Equalf(t, tt.want, digit.ShiftLeft(tt.times).String(), "ShiftLeft(%v)", tt.times)
		})
	}
}

func TestLexoInteger_ShiftRight(t *testing.T) {
	tests := []struct {
		name  string
		digit string
		times int
		want  string
	}{
		{
			name:  "zero times",
			digit: "1234",
			times: 0,
			want:  "1234",
		},
		{
			name:  "one times",
			digit: "1234",
			times: 1,
			want:  "123",
		},
		{
			name:  "more times",
			digit: "1234",
			times: 4,
			want:  "0",
		},
		{
			name:  "negative times",
			digit: "1234",
			times: -2,
			want:  "123400",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			digit, _ := LexoIntegerParse(tt.digit, NewLexoNumeralSystem36())
			assert.Equalf(t, tt.want, digit.ShiftRight(tt.times).String(), "ShiftRight(%v)", tt.times)
		})
	}
}

func TestLexoInteger_Compare(t *testing.T) {
	tests := []struct {
		name  string
		left  string
		right string
		want  int
	}{
		{
			name:  "zero",
			left:  "0",
			right: "0",
			want:  0,
		},
		{
			name:  "equal positive",
			left:  "1234",
			right: "1234",
			want:  0,
		},
		{
			name:  "equal negative",
			left:  "-1234",
			right: "-1234",
			want:  0,
		},
		{
			name:  "positive long left",
			left:  "1234",
			right: "123",
			want:  1,
		},
		{
			name:  "positive long right",
			left:  "123",
			right: "1234",
			want:  -1,
		},
		{
			name:  "positive big left",
			left:  "4321",
			right: "1234",
			want:  1,
		},
		{
			name:  "positive big right",
			left:  "1234",
			right: "4321",
			want:  -1,
		},
		{
			name:  "negative long left",
			left:  "-1234",
			right: "-123",
			want:  -1,
		},
		{
			name:  "negative long right",
			left:  "-123",
			right: "-1234",
			want:  1,
		},
		{
			name:  "negative big left",
			left:  "-1234",
			right: "-4321",
			want:  1,
		},
		{
			name:  "negative big right",
			left:  "-4321",
			right: "-1234",
			want:  -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, _ := LexoIntegerParse(tt.left, NewLexoNumeralSystem36())
			right, _ := LexoIntegerParse(tt.right, NewLexoNumeralSystem36())
			assert.Equalf(t, tt.want, left.Compare(right), "Compare(%v)", tt.right)
		})
	}
}

func TestLexoInteger_Add(t *testing.T) {
	tests := []struct {
		name  string
		left  string
		right string
		want  string
	}{
		{
			name:  "0+0",
			left:  "0",
			right: "0",
			want:  "0",
		},
		{
			name:  "0+1",
			left:  "0",
			right: "1",
			want:  "1",
		},
		{
			name:  "1+0",
			left:  "1",
			right: "0",
			want:  "1",
		},
		{
			name:  "1234+5678",
			left:  "1234",
			right: "5678",
			want:  "68ac",
		},
		{
			name:  "zzzz+zzzz",
			left:  "zzzz",
			right: "zzzz",
			want:  "1zzzy",
		},
		{
			name:  "-2+-3",
			left:  "-2",
			right: "-3",
			want:  "-5",
		},
		{
			name:  "-123+456",
			left:  "-123",
			right: "456",
			want:  "333",
		},
		{
			name:  "123+-456",
			left:  "123",
			right: "-456",
			want:  "-333",
		},
	}
	system := NewLexoNumeralSystem36()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, _ := LexoIntegerParse(tt.left, system)
			right, _ := LexoIntegerParse(tt.right, system)
			res, _ := left.Add(right)
			assert.Equalf(t, tt.want, res.String(), "Add(%v)", tt.right)
		})
	}
}

func TestLexoInteger_Sub(t *testing.T) {
	tests := []struct {
		name  string
		left  string
		right string
		want  string
	}{
		{
			name:  "0-0",
			left:  "0",
			right: "0",
			want:  "0",
		},
		{
			name:  "0-1",
			left:  "0",
			right: "1",
			want:  "-1",
		},
		{
			name:  "1-0",
			left:  "1",
			right: "0",
			want:  "1",
		},
		{
			name:  "1-1",
			left:  "1",
			right: "1",
			want:  "0",
		},
		{
			name:  "1-2",
			left:  "1",
			right: "2",
			want:  "-1",
		},
		{
			name:  "2-1",
			left:  "2",
			right: "1",
			want:  "1",
		},
		{
			name:  "123-1",
			left:  "123",
			right: "1",
			want:  "122",
		},
		{
			name:  "1-123",
			left:  "1",
			right: "123",
			want:  "-122",
		},
		{
			name:  "68ac-1234",
			left:  "68ac",
			right: "1234",
			want:  "5678",
		},
		{
			name:  "1zzzy-zzzz",
			left:  "1zzzy",
			right: "zzzz",
			want:  "zzzz",
		},
		{
			name:  "-2--3",
			left:  "-2",
			right: "-3",
			want:  "1",
		},
		{
			name:  "123--456",
			left:  "123",
			right: "-456",
			want:  "579",
		},
		{
			name:  "-123-456",
			left:  "-123",
			right: "456",
			want:  "-579",
		},
	}
	system := NewLexoNumeralSystem36()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, _ := LexoIntegerParse(tt.left, system)
			right, _ := LexoIntegerParse(tt.right, system)
			res, _ := left.Sub(right)
			assert.Equalf(t, tt.want, res.String(), "Sub(%v)", tt.right)
		})
	}
}

func TestLexoInteger_Multiply(t *testing.T) {
	tests := []struct {
		name  string
		left  string
		right string
		want  string
	}{
		{
			name:  "left zero",
			left:  "0",
			right: "1234",
			want:  "0",
		},
		{
			name:  "right zero",
			left:  "1234",
			right: "0",
			want:  "0",
		},
		{
			name:  "left one",
			left:  "1",
			right: "1234",
			want:  "1234",
		},
		{
			name:  "right one",
			left:  "1234",
			right: "1",
			want:  "1234",
		},
		{
			name:  "positive*positive",
			left:  "10",
			right: "10",
			want:  "100",
		},
		{
			name:  "negative*negative",
			left:  "-1234",
			right: "-5678",
			want:  "5gzpqgw",
		},
		{
			name:  "maxdigit*halfdigit",
			left:  "z",
			right: "i",
			want:  "hi",
		},
	}
	system := NewLexoNumeralSystem36()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, _ := LexoIntegerParse(tt.left, system)
			right, _ := LexoIntegerParse(tt.right, system)
			res, _ := left.Multiply(right)
			assert.Equalf(t, tt.want, res.String(), "Multiply(%v)", right)
		})
	}
}
