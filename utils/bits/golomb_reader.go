package bits

import (
	"fmt"
	"io"
)

// NapToRbsp ...
func NapToRbsp(nal []byte) []byte {
	rbsp := make([]byte, 0, len(nal))
	zeroCount := 0
	for _, b := range nal {
		if zeroCount == 2 && b == 0x3 {
			zeroCount = 0
			continue
		}
		rbsp = append(rbsp, b)
		if b == 0 {
			zeroCount++
		} else {
			zeroCount = 0
		}
	}
	return rbsp
}

type GolombBitReader struct {
	R      io.Reader
	buf    [1]byte
	left   byte
	pos    int
	Debug  bool
	Debug2 bool
}

func (self *GolombBitReader) Pos() int {
	return self.pos
}

func (self *GolombBitReader) Left() byte {
	return self.left
}

func (self *GolombBitReader) CurByte() byte {
	return self.buf[0]
}

func (self *GolombBitReader) ReadBit() (res uint, err error) {
	if self.left == 0 {
		if _, err = self.R.Read(self.buf[:]); err != nil {
			return
		}
		if self.Debug {
			fmt.Printf("got %2x at pos %d\n", self.buf[0], self.pos)
		}
		self.pos++
		self.left = 8
	}
	self.left--
	res = uint(self.buf[0]>>self.left) & 1
	return
}

func (self *GolombBitReader) ReadBits(n int) (res uint, err error) {
	for i := 0; i < n; i++ {
		var bit uint
		if bit, err = self.ReadBit(); err != nil {
			return
		}
		if self.Debug2 {
			fmt.Printf("i %2d ind %2d bit %d ", i, n-i-1, bit)

		}
		res |= bit << uint(n-i-1)
	}
	if self.Debug2 {
		fmt.Printf(" res = %x\n", res)
	}
	return
}

func (self *GolombBitReader) ReadExponentialGolombCode() (res uint, err error) {
	i := 0
	for {
		var bit uint
		if bit, err = self.ReadBit(); err != nil {
			return
		}
		if !(bit == 0 && i < 32) {
			break
		}
		i++
	}
	if res, err = self.ReadBits(i); err != nil {
		return
	}
	res += (1 << uint(i)) - 1
	return
}

func (self *GolombBitReader) ReadSE() (res uint, err error) {
	if res, err = self.ReadExponentialGolombCode(); err != nil {
		return
	}
	if res&0x01 != 0 {
		res = (res + 1) / 2
	} else {
		res = -res / 2
	}
	return
}
