package lexorank

import (
	"errors"
	"fmt"
	"strings"
)

var (
	LexoRankSystem = NewLexoNumeralSystem36()

	zeroDecimal, _  = LexoDecimalParse("0", LexoRankSystem)
	oneDecimal, _   = LexoDecimalParse("1", LexoRankSystem)
	eightDecimal, _ = LexoDecimalParse("8", LexoRankSystem)
	megaDecimal, _  = LexoDecimalParse("1000000", LexoRankSystem)

	minDecimal = zeroDecimal
	maxDecimal = megaDecimal.Sub(oneDecimal)

	MinLexoRank    = NewLexoRank(LexoRankBucket0, minDecimal)
	MaxLexoRank    = NewLexoRank(LexoRankBucket0, maxDecimal)
	MidLexoRank, _ = MinLexoRank.Between(MaxLexoRank)

	initialMinDecimal, _ = LexoDecimalParse("100000", LexoRankSystem)
	initialMaxDecimal, _ = LexoDecimalParse(string(LexoRankSystem.Char(LexoRankSystem.GetBase()-byte(2)))+"00000", LexoRankSystem)
)

type LexoRank struct {
	value   string
	bucket  *LexoRankBucket
	decimal *LexoDecimal
}

func NewLexoRank(bucket *LexoRankBucket, decimal *LexoDecimal) *LexoRank {
	return &LexoRank{
		value:   bucket.String() + "|" + formatDecimal(decimal),
		bucket:  bucket,
		decimal: decimal,
	}
}

func formatDecimal(decimal *LexoDecimal) string {
	formatVal := decimal.String()
	partialIndex := strings.Index(formatVal, string(LexoRankSystem.GetRadixPointChar()))
	if partialIndex < 0 {
		partialIndex = len(formatVal)
		formatVal += string(LexoRankSystem.GetRadixPointChar())
	}
	return strings.Repeat("0", 6-partialIndex) + formatVal
}

func LexoRankParse(str string) (*LexoRank, error) {
	split := strings.Split(str, "|")
	if len(split) != 2 {
		return nil, errors.New("parts not two")
	}
	bucket, err := NewLexoRankBucket(split[0])
	if err != nil {
		return nil, fmt.Errorf("lexo rank bucket: %w", err)
	}
	decimal, err := LexoDecimalParse(split[1], LexoRankSystem)
	if err != nil {
		return nil, fmt.Errorf("lexo decimal parse: %w", err)
	}
	return NewLexoRank(bucket, decimal), nil
}

func (i *LexoRank) Between(other *LexoRank) (*LexoRank, error) {
	if !i.bucket.Equals(other.bucket) {
		return nil, errors.New("between works only within the same bucket")
	}
	cmp := i.decimal.Compare(other.decimal)
	switch {
	case cmp > 0:
		between, err := other.decimal.Between(i.decimal)
		if err != nil {
			return nil, fmt.Errorf("lexo rank between: %w", err)
		}
		return NewLexoRank(i.bucket, between), nil
	case cmp < 0:
		between, err := i.decimal.Between(other.decimal)
		if err != nil {
			return nil, fmt.Errorf("lexo rank between: %w", err)
		}
		return NewLexoRank(i.bucket, between), nil
	}
	return nil, fmt.Errorf("try to rank between issues with same rank this=%s other=%s", i.String(), other.String())
}

func (i *LexoRank) Prev() (*LexoRank, error) {
	if i.IsMax() {
		return NewLexoRank(i.bucket, initialMaxDecimal), nil
	}
	ceilInteger := i.decimal.Ceil()
	ceilDecimal := LexoDecimalMake(ceilInteger, 0)
	nextDecimal := ceilDecimal.Sub(eightDecimal)
	if nextDecimal.Compare(minDecimal) <= 0 {
		nextDecimal, _ = i.decimal.Between(minDecimal)
	}
	return NewLexoRank(i.bucket, nextDecimal), nil
}

func (i *LexoRank) Next() (*LexoRank, error) {
	if i.IsMin() {
		return NewLexoRank(i.bucket, initialMinDecimal), nil
	}
	ceilInteger := i.decimal.Ceil()
	ceilDecimal := LexoDecimalMake(ceilInteger, 0)
	nextDecimal := ceilDecimal.Add(eightDecimal)
	if nextDecimal.Compare(maxDecimal) >= 0 {
		nextDecimal, _ = i.decimal.Between(maxDecimal)
	}
	return NewLexoRank(i.bucket, nextDecimal), nil
}

func (i *LexoRank) String() string {
	return i.value
}

func (i *LexoRank) IsMin() bool {
	return i.decimal.Equals(minDecimal)
}

func (i *LexoRank) IsMax() bool {
	return i.decimal.Equals(maxDecimal)
}
