package lexorank

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexoDecimalParse(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		system  LexoNumeralSystem
		want    *LexoDecimal
		wantErr bool
	}{
		{
			name:   "empty",
			str:    "",
			system: NewLexoNumeralSystem36(),
			want: &LexoDecimal{
				mag: &LexoInteger{
					sys:  NewLexoNumeralSystem36(),
					sign: 0,
					mag:  []byte{0},
				},
				scale: 0,
			},
			wantErr: false,
		},
		{
			name:   "zero",
			str:    "0",
			system: NewLexoNumeralSystem36(),
			want: &LexoDecimal{
				mag: &LexoInteger{
					sys:  NewLexoNumeralSystem36(),
					sign: 0,
					mag:  []byte{0},
				},
				scale: 0,
			},
			wantErr: false,
		},
		{
			name:   "one",
			str:    "1",
			system: NewLexoNumeralSystem36(),
			want: &LexoDecimal{
				mag: &LexoInteger{
					sys:  NewLexoNumeralSystem36(),
					sign: 1,
					mag:  []byte{1},
				},
				scale: 0,
			},
			wantErr: false,
		},
		{
			name:   "pi",
			str:    "3:14159",
			system: NewLexoNumeralSystem36(),
			want: &LexoDecimal{
				mag: &LexoInteger{
					sys:  NewLexoNumeralSystem36(),
					sign: 1,
					mag:  []byte{9, 5, 1, 4, 1, 3},
				},
				scale: 5,
			},
			wantErr: false,
		},
		{
			name:   "-pi",
			str:    "-3:14159",
			system: NewLexoNumeralSystem36(),
			want: &LexoDecimal{
				mag: &LexoInteger{
					sys:  NewLexoNumeralSystem36(),
					sign: -1,
					mag:  []byte{9, 5, 1, 4, 1, 3},
				},
				scale: 5,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LexoDecimalParse(tt.str, tt.system)
			if (err != nil) != tt.wantErr {
				t.Errorf(fmt.Sprintf("LexoDecimalParse(%v, %v) error = %v, wantErr %v", tt.str, tt.system, err, tt.wantErr))
				return
			}
			assert.Equalf(t, tt.want, got, "LexoDecimalParse(%v, %v)", tt.str, tt.system)
		})
	}
}

func TestLexoDecimal_String(t *testing.T) {
	tests := []struct {
		name    string
		decimal *LexoDecimal
		want    string
	}{
		{
			name: "zero",
			decimal: &LexoDecimal{
				mag: &LexoInteger{
					sys:  NewLexoNumeralSystem36(),
					sign: 0,
					mag:  []byte{0},
				},
				scale: 0,
			},
			want: "0",
		},
		{
			name: "pi",
			decimal: &LexoDecimal{
				mag: &LexoInteger{
					sys:  NewLexoNumeralSystem36(),
					sign: 1,
					mag:  []byte{9, 5, 1, 4, 1, 3},
				},
				scale: 5,
			},
			want: "3:14159",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.decimal.String(), "String()")
		})
	}
}

func TestLexoDecimal_SetScale(t *testing.T) {
	tests := []struct {
		name    string
		decimal *LexoDecimal
		scale   int
		want    *LexoDecimal
	}{
		{
			name: "zero scale two",
			decimal: &LexoDecimal{
				mag: &LexoInteger{
					sys:  NewLexoNumeralSystem36(),
					sign: 0,
					mag:  []byte{0},
				},
				scale: 0,
			},
			scale: 2,
			want: &LexoDecimal{
				mag: &LexoInteger{
					sys:  NewLexoNumeralSystem36(),
					sign: 0,
					mag:  []byte{0},
				},
				scale: 0,
			},
		},
		{
			name: "scale more to one",
			decimal: &LexoDecimal{
				mag: &LexoInteger{
					sys:  NewLexoNumeralSystem36(),
					sign: 1,
					mag:  []byte{6, 5, 4, 3, 2, 1},
				},
				scale: 3,
			},
			scale: 1,
			want: &LexoDecimal{
				mag: &LexoInteger{
					sys:  NewLexoNumeralSystem36(),
					sign: 1,
					mag:  []byte{4, 3, 2, 1},
				},
				scale: 1,
			},
		},
		{
			name: "scale more",
			decimal: &LexoDecimal{
				mag: &LexoInteger{
					sys:  NewLexoNumeralSystem36(),
					sign: 1,
					mag:  []byte{6, 5, 4, 3, 2, 1},
				},
				scale: 3,
			},
			scale: 7,
			want: &LexoDecimal{
				mag: &LexoInteger{
					sys:  NewLexoNumeralSystem36(),
					sign: 1,
					mag:  []byte{6, 5, 4, 3, 2, 1},
				},
				scale: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.decimal.SetScale(tt.scale), "SetScale(%v)", tt.scale)
		})
	}
}

func TestLexoDecimal_Multiply(t *testing.T) {
	tests := []struct {
		name  string
		left  string
		right string
		want  string
	}{
		{
			name:  "2*pi",
			left:  "2",
			right: "3:14159",
			want:  "6:282ai",
		},
		{
			name:  "pi*pi",
			left:  "3:14159",
			right: "3:14159",
			want:  "9:6pfe1pb9k9",
		},
	}
	system := NewLexoNumeralSystem36()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, _ := LexoDecimalParse(tt.left, system)
			right, _ := LexoDecimalParse(tt.right, system)
			assert.Equalf(t, tt.want, left.Multiply(right).String(), "Multiply(%v)", right)
		})
	}
}
