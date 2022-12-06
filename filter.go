package bloomzinho

import (
	b "math/bits"
)

type filterParams struct {
	nhsh int //number of hashes
	len  int //number of bits
	bph  int //bits per hashes (indexes)
	ibf  int //iterations before shuffling
}

type bucket_t uint64

type Filter struct {
	state []bucket_t

	filterParams
}

// bpb is bits per bucket (Filter.state element)
const bpb = 64

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
		state: make([]bucket_t, bits/bpb),
		filterParams: filterParams{
			nhsh: hashes,
			bph:  needs,
			ibf:  iters,
			len:  bits,
		},
	}
}

func (f *Filter) AddString(s string) {
	hash := hashS(s)
	h := f.hashToIndexes(hash)

	f.set(h)
}

func (f *Filter) LookupString(s string) bool {
	hash := hashS(s)
	h := f.hashToIndexes(hash)

	return f.lookup(h)
}

func (f *Filter) AddBytes(b []byte) {
	hash := hashB(b)
	h := f.hashToIndexes(hash)

	f.set(h)
}

func (f *Filter) LookupBytes(b []byte) bool {
	hash := hashB(b)
	h := f.hashToIndexes(hash)

	return f.lookup(h)
}

// hashToIndexes as the name says, gets a 64bit hash and turns them into bitfield indexes
//
// all indexes are derived from the same 64 bits
// when all bits have been consumed we just shuffle them and keep going...
// this seem good enough, as long as TestBigHashShuffle passes, it should be fine
// unless the number of items is so big (4 billion?) that they start coliding
// on the original 64bit hash
func (f *Filter) hashToIndexes(hash uint64) []uint {
	idx := make([]uint, 0, 8)
	f.hashTransform(hash, func(u uint) {
		idx = append(idx, u)
	})

	return idx
}

func (f *Filter) hashTransform(hash uint64, cb func(uint)) {
	max := uint64(f.len)
	mask := (uint64(1) << f.bph) - 1

	for i := 0; i < f.nhsh; i++ {
		if i != 0 && i%f.ibf == 0 {
			//shuffles the bits everytime we consume everything
			hash *= 0xc67c_2dcb_6b04_ebb5 //got this from /dev/urandom
		}

		h := (hash & mask) % max
		cb(uint(h))

		//rotate to get the next bits
		//we rotate instead of shifting because we want to keep the bits fo shuffling later
		hash = b.RotateLeft64(hash, f.bph)
	}
}

func (f *Filter) set(idx []uint) {
	for _, i := range idx {
		bucket := i / bpb
		shift := i % bpb

		mask := bucket_t(1 << shift)

		f.state[bucket] |= mask
	}
}

func (f *Filter) lookup(idx []uint) bool {
	for _, v := range idx {
		bucket := v / bpb
		shift := v % bpb

		mask := bucket_t(1 << shift)

		set := f.state[bucket]&mask != 0
		if !set {
			return false
		}
	}

	return len(idx) != 0
}
