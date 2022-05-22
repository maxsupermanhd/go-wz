package main

import (
	"encoding/binary"
	"log"
	"os"
)

// type LengthInteger interface {
// 	~int | ~int8 | ~int16 | ~int32 | ~int64
// 	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
// }

// func readBEBytes[L LengthInteger](f *os.File, n L) ([]byte, error) {

func readBytes(f *os.File, n int) ([]byte, error) {
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
func readByte(f *os.File) (byte, bool, error) {
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

func readUBE32(f *os.File) (uint32, error) {
	b, err := readBytes(f, 4)
	return binary.BigEndian.Uint32(b[0:]), err
}

func decode_uint32_t(b uint8, v uint32, n uint) (bool, uint32) {
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
