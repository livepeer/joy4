package ts

import (
	"time"

	"github.com/livepeer/joy4/av"
	"github.com/livepeer/joy4/format/ts/tsio"
)

type Stream struct {
	av.CodecData

	demuxer *Demuxer
	muxer   *Muxer

	pid        uint16
	streamId   uint8
	streamType uint8

	tsw *tsio.TSWriter
	idx int

	iskeyframe   bool
	pts, dts     time.Duration
	ptsTS, dtsTS int64
	data         []byte
	datalen      int
}
