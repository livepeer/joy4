package ts

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/livepeer/joy4/av"
	"github.com/livepeer/joy4/codec/aacparser"
	"github.com/livepeer/joy4/codec/h264parser"
	"github.com/livepeer/joy4/format/ts/tsio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testDemux(t *testing.T, dm *Demuxer, second bool) []av.Packet {
	assert := assert.New(t)
	require := require.New(t)
	var packets []av.Packet

	streams, err := dm.Streams()
	require.Nil(err)
	assert.Len(streams, 2)
	hcd, ok := streams[0].(h264parser.CodecData)
	require.True(ok)
	assert.Equal(tsio.PTS_HZ, int(hcd.TimeScale()))
	assert.Equal(60., hcd.SPSInfo.FPS)
	assert.Equal(uint(120), hcd.SPSInfo.TimeScale)

	acd, ok := streams[1].(aacparser.CodecData)
	require.True(ok)
	assert.Equal(44100, acd.SampleRate())
	assert.Equal(tsio.PTS_HZ, int(acd.TimeScale()))

	numVFrames := 0
	pkt, err := dm.ReadPacket()
	packets = append(packets, pkt)
	require.Nil(err)
	if pkt.Idx == int8(0) {
		numVFrames++
	}
	assert.True(pkt.IsKeyFrame)
	assert.Equal(int8(0), pkt.Idx)
	assert.Equal(1400*time.Millisecond, pkt.Time)
	assert.Equal(int64(1400*tsio.PTS_HZ/1000), pkt.TimeTS)
	assert.Equal(int64(tsio.PTS_HZ), pkt.TimeScale)
	assert.Equal(33*time.Millisecond, pkt.CompositionTime)
	assert.Equal(int64(33*tsio.PTS_HZ/1000), pkt.CompositionTimeTS)
	pkt, err = dm.ReadPacket()
	packets = append(packets, pkt)
	require.Nil(err)
	if pkt.Idx == int8(0) {
		numVFrames++
	}
	assert.False(pkt.IsKeyFrame)
	assert.Equal(int8(0), pkt.Idx)
	assert.Equal(1416*time.Millisecond, pkt.Time)
	assert.Equal(int64(1416*tsio.PTS_HZ/1000), pkt.TimeTS)
	assert.Equal(int64(tsio.PTS_HZ), pkt.TimeScale)
	assert.Equal(83*time.Millisecond, pkt.CompositionTime)
	assert.Equal(int64(83*tsio.PTS_HZ/1000), pkt.CompositionTimeTS)
	for {
		pkt, err = dm.ReadPacket()
		packets = append(packets, pkt)
		require.Nil(err)
		if pkt.Idx == int8(0) {
			numVFrames++
		}
		if pkt.Idx > 0 {
			break
		}
	}
	assert.Equal(int8(1), pkt.Idx)
	assert.Equal(1400*time.Millisecond, pkt.Time)
	assert.Equal(int64(1400*tsio.PTS_HZ/1000), pkt.TimeTS)
	assert.Equal(0*time.Millisecond, pkt.CompositionTime)
	pkt, err = dm.ReadPacket()
	packets = append(packets, pkt)
	require.Nil(err)
	if pkt.Idx == int8(0) {
		numVFrames++
	}
	assert.Equal(int8(1), pkt.Idx)
	if !second {
		assert.Equal(1423219954*time.Nanosecond, pkt.Time)
	}
	assert.Equal(int64(1400*tsio.PTS_HZ/1000)+int64(1024*tsio.PTS_HZ/acd.SampleRate()), pkt.TimeTS)
	for {
		pkt, err = dm.ReadPacket()
		if err == io.EOF {
			break
		}
		packets = append(packets, pkt)
		require.Nil(err)
		if pkt.Idx == int8(0) {
			numVFrames++
		}
	}
	assert.Equal(5, numVFrames)
	return packets
}

func TestTs(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	cwd, err := os.Getwd()
	require.Nil(err)
	fmt.Printf("cwd: %s\n", cwd)
	file, err := os.Open("../../data/short0.ts")
	require.Nil(err)

	dm := NewDemuxer(file)
	// dump(dm)
	// file2, err := os.Open("../../data/short0.ts")
	// require.Nil(err)
	// dm = NewDemuxer(file2)

	packets := testDemux(t, dm, false)
	assert.Equal(10, len(packets))
	// now mux back
	buf := new(bytes.Buffer)
	mx := NewMuxer(buf)
	streams, err := dm.Streams()
	require.Nil(err)
	err = mx.WriteHeader(streams)
	require.Nil(err)
	for _, pkt := range packets {
		err = mx.WritePacket(pkt)
		require.Nil(err)
	}
	mx.WriteTrailer()
	require.Nil(err)

	// fmt.Println("---> second time")
	dm2 := NewDemuxer(buf)
	// dump(dm2)
	// assert.True(false)
	packets = testDemux(t, dm2, true)
	assert.Equal(10, len(packets))
}

func dump(dm *Demuxer) {
	num := 0
	for {
		pkt, err := dm.ReadPacket()
		if err == io.EOF {
			break
		}
		fmt.Printf("Packet idx %d key %v time %s cts %s (%s) raw time %d cts %d timeScale %d\n",
			pkt.Idx, pkt.IsKeyFrame, pkt.Time, pkt.CompositionTime, pkt.Time+pkt.CompositionTime, pkt.TimeTS,
			pkt.CompositionTimeTS, pkt.TimeScale)
		num++
	}
	fmt.Printf("Got %d packets\n", num)
}
