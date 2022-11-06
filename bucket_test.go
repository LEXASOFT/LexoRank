package lexorank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexoRankBucket_String(t *testing.T) {
	tests := []struct {
		name   string
		bucket *LexoRankBucket
		want   string
	}{
		{
			name:   "bucket_0",
			bucket: LexoRankBucket0,
			want:   "0",
		},
		{
			name:   "bucket_1",
			bucket: LexoRankBucket1,
			want:   "1",
		},
		{
			name:   "bucket_2",
			bucket: LexoRankBucket2,
			want:   "2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.bucket.String(), "String()")
		})
	}
}
