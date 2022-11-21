package bloomzinho

import (
	"hash/maphash"
	b "math/bits"
)

type Filter struct {
	seed  maphash.Seed
	state []uint8
	nhsh  int
	bph   int
}

func NewFilter(bits, hashes int) *Filter {
	if bits%8 != 0 {
		bits += 8
	}
	if bits < 1 {
		panic("bits need to be a higher than 0")
	}
	lz := b.LeadingZeros64(uint64(bits))
	needs := 64 - lz
	if needs*hashes > 64 {
		//TODO do 2 hashes?
		panic("too... much... data...")
	}
	return &Filter{
		seed:  maphash.MakeSeed(),
		state: make([]uint8, bits/8),
		nhsh:  hashes,
		bph:   needs,
	}
}

func (f *Filter) AddString(s string) {
	h := f.doHashString(s)

	for _, v := range h {
		f.set(v)
	}
}

func (f *Filter) LookupString(s string) bool {
	h := f.doHashString(s)

	for _, v := range h {
		if !f.isSet(v) {
			return false
		}
	}
	return true
}

func (f *Filter) doHashString(s string) []int {
	hash := maphash.String(f.seed, s)

	bph := f.bph
	max := uint64(len(f.state) * 8)
	mask := (uint64(1) << f.bph) - 1

	var ret []int

	for i := 0; i < f.nhsh; i++ {
		h := hash & mask
		h %= max
		ret = append(ret, int(h))

		hash >>= bph
	}

	return ret
}

func (f *Filter) set(i int) {
	bucket := i / 8
	shift := i % 8

	mask := uint8(1 << shift)

	f.state[bucket] |= mask
}

func (f *Filter) isSet(i int) bool {
	bucket := i / 8
	shift := i % 8

	mask := uint8(1 << shift)

	set := f.state[bucket]&mask != 0

	return set
}
