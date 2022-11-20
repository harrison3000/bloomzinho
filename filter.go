package bloomzinho

import (
	"encoding/binary"
	"hash"
	"hash/fnv"
)

type Filter struct {
	h     hash.Hash
	state []uint8
	nh    int
}

func NewFilter(bits, hashes int) *Filter {
	if bits%8 != 0 {
		bits += 8
	}
	return &Filter{
		h:     fnv.New128a(),
		state: make([]uint8, bits/8),
		nh:    hashes,
	}
}

func (f *Filter) AddString(s string) {
	h := f.doHashString(s)

	for _, v := range h {
		f.set(v)
	}
}

func (f *Filter) LookupString(s string) bool {
	return false
}

func (f *Filter) doHashString(s string) []int {
	f.h.Reset()
	f.h.Write([]byte(s))
	hash := f.h.Sum(nil)

	nbits := len(f.state) * 8

	if nbits > 65000 || f.nh > 8 {
		panic("we dont do that here (yet)")
	}

	var ret []int

	for i := 0; i < f.nh; i++ {
		idx := i * 2

		h := int(binary.LittleEndian.Uint16(hash[idx:]))

		h %= nbits

		ret = append(ret, h)
	}

	return ret
}

func (f *Filter) set(i int) {
	bucket := i / 8
	shift := i % 8

	mask := uint8(1 << shift)

	f.state[bucket] |= mask
}
