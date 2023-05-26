package wznet

import (
	"encoding/binary"
	"io"
)

func ReadBytes(f io.Reader, n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := io.ReadFull(f, b)
	return b, err
}
func ReadByte(f io.Reader) (byte, bool, error) {
	var p [1]byte
	_, err := f.Read(p[:])
	if err != nil {
		return p[0], false, err
	}
	return p[0], true, nil
}

func ReadUBE32(f io.Reader) (ret uint32, err error) {
	err = binary.Read(f, binary.BigEndian, &ret)
	return
}

var (
	table_uint32_t_a = []uint32{78, 95, 32, 70, 0}
	table_uint32_t_m = []uint32{1, 78, 7410, 237120, 16598400}
)

func Decode_uint32_t(b uint8, v uint32, n uint) (bool, uint32) {
	a := table_uint32_t_a[n]
	m := table_uint32_t_m[n]
	isLastByte := uint32(b) < 256-a
	if isLastByte {
		v += uint32(b) * m
	} else {
		v += (256 - a + 255 - uint32(b)) * m
	}
	return isLastByte, v
}

func NETreadU8(r io.Reader) (ret uint8, err error) {
	err = binary.Read(r, binary.BigEndian, &ret)
	return
}

func NETreadU16(r io.Reader) (ret uint16, err error) {
	err = binary.Read(r, binary.BigEndian, &ret)
	return
}

func NETreadU32(r io.Reader) (ret uint32, err error) {
	end := false
	for n := uint(0); !end; n++ {
		b := byte(0)
		err = binary.Read(r, binary.BigEndian, &b)
		if err != nil {
			return 0, err
		}
		end, ret = Decode_uint32_t(b, ret, n)
	}
	return
}

func NETreadS32(r io.Reader) (ret int32, err error) {
	v, err := NETreadU32(r)
	if err != nil {
		return 0, err
	}
	if v%2 == 0 {
		return int32(v / 2), nil
	} else {
		return -(int32(v/2) - 1), nil
	}
}

func NETstring(r io.Reader) (ret string, err error) {
	ret = ""
	len, err := NETreadU32(r)
	if err != nil {
		return ret, err
	}
	for ; len > 0; len-- {
		c, err := NETreadU16(r)
		if err != nil {
			return ret, err
		}
		ret += string(rune(c))
	}
	return ret, nil
}
