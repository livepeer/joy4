package h264parser

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	assert := assert.New(t)
	var typ int
	var nalus [][]byte

	annexbFrame, _ := hex.DecodeString("00000001223322330000000122332233223300000133000001000001")
	nalus, typ = SplitNALUs(annexbFrame)
	assert.Len(nalus, 3)
	assert.Equal(NALU_ANNEXB, typ)

	avccFrame, _ := hex.DecodeString(
		"00000008aabbccaabbccaabb00000001aa",
	)
	nalus, typ = SplitNALUs(avccFrame)
	t.Log(typ, len(nalus))
	assert.Len(nalus, 2)
	assert.Equal(NALU_AVCC, typ)
}
