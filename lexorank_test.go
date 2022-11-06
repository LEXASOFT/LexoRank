package lexorank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexoRankBetween(t *testing.T) {
	minNext, _ := MinLexoRank.Next()
	maxPrev, _ := MaxLexoRank.Prev()

	tests := []struct {
		name  string
		left  string
		right string
		want  string
	}{
		{
			name:  "Min <-> Max",
			left:  MinLexoRank.String(),
			right: MaxLexoRank.String(),
			want:  "0|hzzzzz:",
		},
		{
			name:  "Min <-> Next",
			left:  MinLexoRank.String(),
			right: minNext.String(),
			want:  "0|0i0000:",
		},
		{
			name:  "Max <-> GetPrev",
			left:  MaxLexoRank.String(),
			right: maxPrev.String(),
			want:  "0|yzzzzz:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, _ := LexoRankParse(tt.left)
			right, _ := LexoRankParse(tt.right)
			got, _ := left.Between(right)
			assert.Equalf(t, tt.want, got.String(), "Between(%v)", right)
		})
	}
}

func TestLexoRankParse(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  *LexoRank
	}{
		{
			name:  "0|hzzzzz:",
			value: "0|hzzzzz:",
			want: &LexoRank{
				value:  "0|hzzzzz:",
				bucket: LexoRankBucket0,
				decimal: &LexoDecimal{
					mag: &LexoInteger{
						sys:  LexoRankSystem,
						sign: 1,
						mag:  []byte{35, 35, 35, 35, 35, 17},
					},
					scale: 0,
				},
			},
		},
		{
			name:  "0|i00000:",
			value: "0|i00000:",
			want: &LexoRank{
				value:  "0|i00000:",
				bucket: LexoRankBucket0,
				decimal: &LexoDecimal{
					mag: &LexoInteger{
						sys:  LexoRankSystem,
						sign: 1,
						mag:  []byte{0, 0, 0, 0, 0, 18},
					},
					scale: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := LexoRankParse(tt.value)
			assert.Equalf(t, tt.want, got, "LexoRankParse(%v)", tt.value)
		})
	}
}

func TestLexoRank_Between(t *testing.T) {
	tests := []struct {
		name  string
		left  string
		right string
		want  string
	}{
		{
			name:  "integer",
			left:  "0|hzzzzz:",
			right: "0|i0000f:",
			want:  "0|i00007:",
		},
		{
			name:  "decimal",
			left:  "0|i00001:",
			right: "0|i00002:",
			want:  "0|i00001:i",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, _ := LexoRankParse(tt.left)
			right, _ := LexoRankParse(tt.right)
			got, _ := left.Between(right)
			assert.Equalf(t, tt.want, got.String(), "Between(%v)", tt.right)
		})
	}
}

func TestLexoRank_Next(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "0|hzzzzz:",
			value: "0|hzzzzz:",
			want:  "0|i00007:",
		},
		{
			name:  "0|i00007:",
			value: "0|i00007:",
			want:  "0|i0000f:",
		},
		{
			name:  "0|i0000f:",
			value: "0|i0000f:",
			want:  "0|i0000n:",
		},
		{
			name:  "0|i0000n:",
			value: "0|i0000n:",
			want:  "0|i0000v:",
		},
		{
			name:  "0|i0000v:",
			value: "0|i0000v:",
			want:  "0|i00013:",
		},
		{
			name:  "0|i00013:",
			value: "0|i00013:",
			want:  "0|i0001b:",
		},
		{
			name:  "0|i0001b:",
			value: "0|i0001b:",
			want:  "0|i0001j:",
		},
		{
			name:  "0|i0001j:",
			value: "0|i0001j:",
			want:  "0|i0001r:",
		},
		{
			name:  "0|i0001r:",
			value: "0|i0001r:",
			want:  "0|i0001z:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, _ := LexoRankParse(tt.value)
			got, _ := i.Next()
			assert.Equalf(t, tt.want, got.String(), "Next()")
		})
	}
}
