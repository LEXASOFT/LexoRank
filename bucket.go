package lexorank

import "fmt"

var (
	LexoRankBucket0, _ = NewLexoRankBucket("0")
	LexoRankBucket1, _ = NewLexoRankBucket("1")
	LexoRankBucket2, _ = NewLexoRankBucket("2")

	LexoRankBuckets = []*LexoRankBucket{LexoRankBucket0, LexoRankBucket1, LexoRankBucket2}
)

type LexoRankBucket struct {
	value *LexoInteger
}

func NewLexoRankBucket(str string) (*LexoRankBucket, error) {
	value, err := LexoIntegerParse(str, LexoRankSystem)
	if err != nil {
		return nil, fmt.Errorf("parse int error: %w", err)
	}
	return &LexoRankBucket{value: value}, nil
}

func (b *LexoRankBucket) String() string {
	return b.value.String()
}

func (b *LexoRankBucket) Equals(other *LexoRankBucket) bool {
	if b == other {
		return true
	}
	if other == nil {
		return false
	}
	return b.value.Equals(other.value)
}
