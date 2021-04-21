package pio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVec(t *testing.T) {
	assert := assert.New(t)
	vec := [][]byte{{1, 2, 3}, {4, 5, 6, 7, 8, 9}, {10, 11, 12, 13}}
	assert.Equal(13, VecLen(vec))

	vec = VecSlice(vec, 2, -1)
	assert.Equal([][]byte{{3}, {4, 5, 6, 7, 8, 9}, {10, 11, 12, 13}}, vec)

	vec = VecSlice(vec, 2, 10)
	assert.Equal([][]byte{{5, 6, 7, 8, 9}, {10, 11, 12}}, vec)

	vec = VecSlice(vec, 8, 8)
	assert.Len(vec, 0)
}
