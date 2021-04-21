// +build !ffmpeg

package ffmpeg

import (
	"errors"
	"image"

	"github.com/livepeer/joy4/av"
)

type VideoDecoder struct {
	Extradata []byte
}

type VideoFrame struct {
	Image image.YCbCr
}

func NewVideoDecoder(stream av.CodecData) (dec *VideoDecoder, err error) {
	return &VideoDecoder{}, nil
}

func (self *VideoDecoder) Decode(pkt []byte) (img *VideoFrame, err error) {
	return nil, errors.New("not implemented")
}
