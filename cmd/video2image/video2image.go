package main

import (
	"fmt"
	"io"

	"github.com/livepeer/joy4/av"
	"github.com/livepeer/joy4/av/avutil"
	"github.com/livepeer/joy4/cgo/ffmpeg"
	"github.com/livepeer/joy4/format"
)

func init() {
	format.RegisterAll()
}

func main() {
	fmt.Println("start")
	file, err := avutil.Open("v108063.ts")
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
	// var dec *ffmpeg.AudioDecoder
	var vdec *ffmpeg.VideoDecoder
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
		if st.Type() == av.H264 {
			vdec, err = ffmpeg.NewVideoDecoder(st)
			if err != nil {
				panic(err)
			}
		}
	}
	fmt.Printf("Video idx %d audio idx %d\n", videoidx, audioidx)
	for {
		pkt, err := file.ReadPacket()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading packet")
			panic(err)
		}
		fmt.Printf("Packet idx %d key %v time %s\n", pkt.Idx, pkt.IsKeyFrame, pkt.Time)
		if pkt.Idx == int8(videoidx) {
			frame, err := vdec.Decode(pkt.Data)
			if err != nil {
				panic(err)
			}
			fmt.Println("decode samples", frame.Image)
		}
	}
}
