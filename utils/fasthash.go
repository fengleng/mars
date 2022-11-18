package utils

import (
	"encoding/binary"
)

func mix(i uint64) uint64 {
	h := i
	h ^= h >> 23
	h *= 0x2127599bf4325c37
	h ^= h >> 47
	return h
}

//FastHash64 FastHash64
func FastHash64(rawData []byte, seed uint64) uint64 {
	m := uint64(0x880355f21e6d1965)
	l := len(rawData)
	h := seed ^ (uint64(l) * m)
	var v uint64
	i := 0
	for i < l-8 {
		bb := make([]byte, 0)
		bb = append(bb, rawData[i:i+8]...)
		v = binary.LittleEndian.Uint64(bb)
		h ^= mix(v)
		h *= m
		i += 8
	}
	t := make([]byte, 8)
	ti := 0
	for k := i; k < l; k++ {
		t[ti] = rawData[k]
		ti++
	}
	v = binary.LittleEndian.Uint64(t)
	h ^= mix(v)
	h *= m
	return mix(h)
}

//FastHash32 FastHash32
func FastHash32(rawData []byte, seed uint32) uint32 {
	h := FastHash64(rawData, uint64(seed))
	return uint32(h - (h >> 32))
}
