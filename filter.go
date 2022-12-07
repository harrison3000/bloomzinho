package bloomzinho

import (
	"errors"
	b "math/bits"
)

type filterParams struct {
	nhsh int //number of hashes
	len  int //number of bits
}

type bucket_t uint64

type Filter struct {
	state []bucket_t

	bph int //bits per hashes (indexes)
	ibf int //iterations before shuffling
	filterParams
}

// bpb is bits per bucket (Filter.state element)
const bpb = 64

func NewFilter(bits, hashes int) (*Filter, error) {
	if bits < 1 || hashes < 1 {
		return nil, errors.New("params needs to be a higher than 0")
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
		bph:   needs,
		ibf:   iters,
		filterParams: filterParams{
			nhsh: hashes,
			len:  bits,
		},
	}, nil
}

func (f *Filter) AddString(s string) {
	hash := hashS(s)

	f.hashTransform(hash, f.set)
}

func (f *Filter) LookupString(s string) bool {
	hash := hashS(s)

	return f.hashTransform(hash, f.lookup)
}

func (f *Filter) AddBytes(b []byte) {
	hash := hashB(b)

	f.hashTransform(hash, f.set)
}

func (f *Filter) LookupBytes(b []byte) bool {
	hash := hashB(b)

	return f.hashTransform(hash, f.lookup)
}

// hashTransform gets a 64bit hash and turns it into bitfield indexes
//
// for each index the callback function is called
// all indexes are derived from the same 64 bits
// when all bits have been consumed we just shuffle them and keep going...
// this seem good enough, as long as TestBigHashShuffle passes, it should be fine
// unless the number of items is so big (4 billion, because of the birthday paradox (?))
// that they start coliding on the original 64bit hash
func (f *Filter) hashTransform(hash uint64, cb func(uint) bool) bool {
	max := uint64(f.len)
	mask := (uint64(1) << f.bph) - 1

	for i := 0; i < f.nhsh; i++ {
		if i != 0 && i%f.ibf == 0 {
			//shuffles the bits everytime we consume everything
			hash *= 0xc67c_2dcb_6b04_ebb5 //got this from /dev/urandom
		}

		h := (hash & mask) % max
		if !cb(uint(h)) {
			return false
		}

		//rotate to get the next bits
		//we rotate instead of shifting because we want to keep the bits fo shuffling later
		//also, rotating is somehow faster than shiffting 🤷
		hash = b.RotateLeft64(hash, f.bph)
	}

	return true
}

func (f *Filter) set(i uint) bool {
	bucket := i / bpb
	shift := i % bpb

	mask := bucket_t(1 << shift)

	f.state[bucket] |= mask

	return true
}

func (f *Filter) lookup(i uint) bool {
	bucket := i / bpb
	shift := i % bpb

	mask := bucket_t(1 << shift)

	set := f.state[bucket]&mask != 0
	return set
}
