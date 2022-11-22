package bloomzinho

import (
	"hash/maphash"
	b "math/bits"
)

type Filter struct {
	seed  maphash.Seed
	state []uint8

	nhsh int //number of hashes
	bph  int //bits per hashes (indexes)
	ibf  int //iterations before shuffling
}

// bpb is bits per bucket (Filter.state element)
const bpb = 8

func NewFilter(bits, hashes int) *Filter {
	if bits < 1 || hashes < 1 {
		panic("params needs to be a higher than 0")
	}

	//round the number of bits to the next multiple of bpd
	if mod := bits % bpb; mod != 0 {
		bits += bpb - mod
	}

	//-1 because arrays are 0 idexed
	//power of 2 number of bits would unnecessarily need an extra bit without this
	needs := b.Len64(uint64(bits) - 1)
	iters := 64 / needs

	return &Filter{
		seed:  maphash.MakeSeed(),
		state: make([]uint8, bits/bpb),
		nhsh:  hashes,
		bph:   needs,
		ibf:   iters,
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

// this function was fine tuned to keep the cost just low enough for inlining
// inlining avoids the return slice escaping to the heap
// TODO explain the idea behind shuffling
// spoiler: something to do with the index orders
func (f *Filter) hashToIndexes(hash uint64) []int {
	max := uint64(len(f.state) * bpb)
	mask := (uint64(1) << f.bph) - 1
	idx := make([]int, 0, 8)

	for i := 0; i < f.nhsh; i++ {
		if i != 0 && i%f.ibf == 0 {
			//shuffles the bits everytime we consume everything
			hash *= 0xc67c_2dcb_6b04_ebb5 //got this from /dev/urandom
		}

		h := (hash & mask) % max
		idx = append(idx, int(h))

		//rotate to get the next bits
		//we rotate instead of shifting because we want to keep the bits fo shuffling later
		hash = b.RotateLeft64(hash, f.bph)
	}

	return idx
}

func (f *Filter) set(i int) {
	bucket := i / bpb
	shift := i % bpb

	mask := uint8(1 << shift)

	f.state[bucket] |= mask
}

func (f *Filter) lookup(idx []int) bool {
	for _, v := range idx {
		bucket := v / bpb
		shift := v % bpb

		mask := uint8(1 << shift)

		set := f.state[bucket]&mask != 0
		if !set {
			return false
		}
	}

	return len(idx) != 0
}
