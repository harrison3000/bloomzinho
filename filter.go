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
	if bits < 1 || hashes < 1 {
		panic("params needs to be a higher than 0")
	}

	if mod := bits % 8; mod != 0 {
		bits += 8 - mod
	}

	//-1 because arrays are 0 idexed
	//power of 2 number of bits would unnecessarily need an extra bit without this
	lz := b.LeadingZeros64(uint64(bits) - 1)
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
	hash := maphash.String(f.seed, s)
	h := f.hashToIndexes(hash)

	for _, v := range h {
		f.set(v)
	}
}

func (f *Filter) LookupString(s string) bool {
	hash := maphash.String(f.seed, s)
	h := f.hashToIndexes(hash)

	return f.lookup(h)
}

func (f *Filter) AddBytes(b []byte) {
	hash := maphash.Bytes(f.seed, b)
	h := f.hashToIndexes(hash)

	for _, v := range h {
		f.set(v)
	}
}

func (f *Filter) LookupBytes(b []byte) bool {
	hash := maphash.Bytes(f.seed, b)
	h := f.hashToIndexes(hash)

	return f.lookup(h)
}

func (f *Filter) hashToIndexes(hash uint64) []int {
	bph := f.bph
	max := uint64(len(f.state) * 8)
	mask := (uint64(1) << f.bph) - 1
	idx := make([]int, 0, 8)

	for i := 0; i < f.nhsh; i++ {
		h := hash & mask
		h %= max
		idx = append(idx, int(h))

		hash >>= bph
	}

	return idx
}

func (f *Filter) set(i int) {
	bucket := i / 8
	shift := i % 8

	mask := uint8(1 << shift)

	f.state[bucket] |= mask
}

func (f *Filter) lookup(idx []int) bool {
	for _, v := range idx {
		bucket := v / 8
		shift := v % 8

		mask := uint8(1 << shift)

		set := f.state[bucket]&mask != 0
		if !set {
			return false
		}
	}

	return len(idx) != 0
}
