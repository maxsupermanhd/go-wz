package wznet

import (
	"encoding/binary"
	"io"
	"log"
)

func ReadBytes(f io.Reader, n int) ([]byte, error) {
	if n <= 0 {
		if n < 0 {
			log.Printf("ReadBEBytes called with negative length! (%d)", n)
		}
		return []byte{}, nil
	}
	b := make([]byte, n)
	r, err := f.Read(b)
	if err != nil {
		return []byte{}, err
	}
	if n != r {
		log.Printf("Read bytes failed to read requested %d, got only %d", n, r)
	}
	return b, err
}
func ReadByte(f io.Reader) (byte, bool, error) {
	b := make([]byte, 1)
	r, err := f.Read(b)
	if err != nil {
		return b[0], false, err
	}
	if r != 1 {
		return b[0], false, err
	}
	return b[0], true, err
}

func ReadUBE32(f io.Reader) (uint32, error) {
	b, err := ReadBytes(f, 4)
	return binary.BigEndian.Uint32(b[0:]), err
}

func Decode_uint32_t(b uint8, v uint32, n uint) (bool, uint32) {
	table_uint32_t_a := []uint32{78, 95, 32, 70, 0}
	table_uint32_t_m := []uint32{1, 78, 7410, 237120, 16598400}
	a := table_uint32_t_a[n]
	m := table_uint32_t_m[n]

	isLastByte := uint32(b) < 256-a
	// log.Printf("Decoding byte %02x (%2d) is last %t", b, n, isLastByte)
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
