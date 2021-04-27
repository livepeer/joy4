package mp4

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/dsnet/golib/memfile"
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
	assert.Equal(44100, int(acd.TimeScale()))

	numVFrames := 0
	pkt, err := dm.ReadPacket()
	packets = append(packets, pkt)
	require.Nil(err)
	assert.Equal(int8(0), pkt.Idx)
	numVFrames++
	assert.True(pkt.IsKeyFrame)
	assert.Equal(int8(0), pkt.Idx)
	assert.Equal(0*time.Millisecond, pkt.Time)
	assert.Equal(int64(0*tsio.PTS_HZ/1000), pkt.TimeTS)
	assert.Equal(int64(tsio.PTS_HZ), pkt.TimeScale)
	assert.Equal(33*time.Millisecond, pkt.CompositionTime)
	assert.Equal(int64(33*tsio.PTS_HZ/1000), pkt.CompositionTimeTS)
	pkt, err = dm.ReadPacket()
	packets = append(packets, pkt)
	require.Nil(err)
	assert.False(pkt.IsKeyFrame)
	assert.Equal(int8(1), pkt.Idx)
	assert.Equal(0*time.Millisecond, pkt.Time)
	assert.Equal(int64(0*tsio.PTS_HZ/1000), pkt.TimeTS)
	assert.Equal(int64(44100), pkt.TimeScale)
	assert.Equal(0*time.Millisecond, pkt.CompositionTime)
	assert.Equal(int64(0*tsio.PTS_HZ/1000), pkt.CompositionTimeTS)

	pkt, err = dm.ReadPacket()
	require.Nil(err)
	packets = append(packets, pkt)
	assert.Equal(int8(0), pkt.Idx)
	numVFrames++
	assert.False(pkt.IsKeyFrame)
	assert.Equal(16*time.Millisecond, pkt.Time)
	assert.Equal(int64(16*tsio.PTS_HZ/1000), pkt.TimeTS)
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
	assert.Equal(int64(44100), pkt.TimeScale)
	assert.Equal(23219954*time.Nanosecond, pkt.Time)
	assert.Equal(int64(1024), pkt.TimeTS)
	assert.Equal(0*time.Millisecond, pkt.CompositionTime)
	pkt, err = dm.ReadPacket()
	packets = append(packets, pkt)
	require.Nil(err)
	assert.Equal(int8(0), pkt.Idx)
	numVFrames++
	// if !second {
	assert.Equal(33*time.Millisecond, pkt.Time)
	// }
	assert.Equal(int64(2970), pkt.TimeTS)
	assert.Equal(int64(2970), pkt.CompositionTimeTS)
	var i int
	for {
		if i == 3 {
			fmt.Println("do")
		}
		pkt, err = dm.ReadPacket()
		if err == io.EOF {
			fmt.Printf("end at i=%d\n", i)
			// pkt, err = dm.ReadPacket()
			// fmt.Printf("On second err=%v raw time %d\n", err, pkt.TimeTS)
			break
		}
		i++
		packets = append(packets, pkt)
		require.Nil(err)
		if pkt.Idx == int8(0) {
			numVFrames++
		}
	}
	assert.Equal(5, numVFrames)
	return packets
}

func TestMp4(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	cwd, err := os.Getwd()
	require.Nil(err)
	fmt.Printf("cwd: %s\n", cwd)
	file, err := os.Open("../../data/short0.mp4")
	require.Nil(err)

	dm := NewDemuxer(file)
	dump(dm)
	file2, err := os.Open("../../data/short0.mp4")
	require.Nil(err)
	dm = NewDemuxer(file2)

	packets := testDemux(t, dm, false)
	require.Equal(10, len(packets))
	// now mux back
	mf := &memfile.File{}
	mx := NewMuxer(mf)
	streams, err := dm.Streams()
	require.Nil(err)
	err = mx.WriteHeader(streams)
	require.Nil(err)
	for i, pkt := range packets {
		err = mx.WritePacket(pkt)
		require.Nil(err)
		fmt.Printf("Wrote packet i=%d idx=%d\n", i, pkt.Idx)
	}
	mx.WriteTrailer()
	require.Nil(err)

	fmt.Println("---> second time")
	mf.Seek(0, 0)
	dm2 := NewDemuxer(mf)
	dump(dm2)
	// assert.True(false)
	mf.Seek(0, 0)
	// ioutil.WriteFile("zf1.mp4", mf.Bytes(), 0644)
	// mf.Seek(0, 0)

	dm2 = NewDemuxer(mf)
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
		fmt.Printf("Packet idx %d key %v time %s cts %s (%s) raw time %d cts %d (%d) timeScale %d\n",
			pkt.Idx, pkt.IsKeyFrame, pkt.Time, pkt.CompositionTime, pkt.Time+pkt.CompositionTime, pkt.TimeTS,
			pkt.CompositionTimeTS, pkt.TimeTS+pkt.CompositionTimeTS, pkt.TimeScale)
		num++
	}
	fmt.Printf("Got %d packets\n", num)
}
