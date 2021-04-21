package main

import (
	"fmt"
	"io"

	"github.com/livepeer/joy4/av"
	"github.com/livepeer/joy4/av/avutil"
	"github.com/livepeer/joy4/format"
)

func init() {
	format.RegisterAll()
}

func main() {
	fmt.Println("start")
	// file, err := avutil.Open("probe.ts")
	file, err := avutil.Open("26.ts")
	fmt.Printf("%T\n", file)
	hd := file.(*avutil.HandlerDemuxer)
	fmt.Printf("%T\n", hd.Demuxer)
	if err != nil {
		panic(err)
	}
	// sc := newSegmentsCounter(segLen, nil, recordSegmentsDurations, nil)
	// filters := pktque.Filters{sc}
	// src := &pktque.FilterDemuxer{Demuxer: file, Filter: filters}
	var streams []av.CodecData
	var videoidx, audioidx int
	if streams, err = file.Streams(); err != nil {
		fmt.Println("Can't count segments in source file")
		panic(err)
	}
	for i, st := range streams {
		if st == nil {
			continue
		}
		fmt.Println(i, " stream is of type ", st.Type())
		if st.Type().IsAudio() {
			audioidx = i
		}
		if st.Type().IsVideo() {
			videoidx = i
		}
		if vc, ok := st.(av.VideoCodecData); ok {
			fmt.Printf("w: %d, h: %d\n", vc.Width(), vc.Height())
		}
	}
	fmt.Printf("Video idx %d audio idx %d\n", videoidx, audioidx)
	for {
		x, err := file.ReadPacket()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading packet")
			panic(err)
		}
		if x.Idx == int8(videoidx) {
			fmt.Printf("Packet idx %d key %v time %s cts %s (%s) raw time %d cts %d\n",
				x.Idx, x.IsKeyFrame, x.Time, x.CompositionTime, x.Time+x.CompositionTime, x.TimeB, x.CompositionTimeB)
		}
	}
}
