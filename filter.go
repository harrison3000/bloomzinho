package bloomzinho

import (
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

}

func (f *Filter) LookupString(s string) bool {
	return false
}
