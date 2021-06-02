// +build !ffmpeg

package ffmpeg

import (
	"errors"

	"github.com/livepeer/joy4/av"
)

type AudioDecoder struct {
	ChannelLayout av.ChannelLayout
	SampleFormat  av.SampleFormat
	SampleRate    int
	Extradata     []byte
}

type AudioEncoder struct {
	SampleRate       int
	Bitrate          int
	ChannelLayout    av.ChannelLayout
	SampleFormat     av.SampleFormat
	FrameSampleCount int
}

func NewAudioDecoder(codec av.AudioCodecData) (dec *AudioDecoder, err error) {
	return &AudioDecoder{}, nil
}

func (self *AudioDecoder) Close() {
}

func (self *AudioDecoder) Decode(pkt []byte) (gotframe bool, frame av.AudioFrame, err error) {
	err = errors.New("not implemented")
	return
}

func NewAudioEncoderByName(name string) (enc *AudioEncoder, err error) {
	return &AudioEncoder{}, nil
}

func (self *AudioEncoder) Close() {
}

func (self *AudioEncoder) CodecData() (codec av.AudioCodecData, err error) {
	err = errors.New("not implemented")
	return
}

func (self *AudioEncoder) Encode(frame av.AudioFrame) (pkts [][]byte, err error) {
	err = errors.New("not implemented")
	return
}

func (self *AudioEncoder) GetOption(key string, val interface{}) (err error) {
	err = errors.New("not implemented")
	return
}

func (self *AudioEncoder) SetBitrate(bitrate int) (err error) {
	return
}

func (self *AudioEncoder) SetChannelLayout(ch av.ChannelLayout) (err error) {
	self.ChannelLayout = ch
	return
}

func (self *AudioEncoder) SetOption(key string, val interface{}) (err error) {
	return
}

func (self *AudioEncoder) SetSampleFormat(fmt av.SampleFormat) (err error) {
	self.SampleFormat = fmt
	return
}

func (self *AudioEncoder) SetSampleRate(rate int) (err error) {
	self.SampleRate = rate
	return
}
