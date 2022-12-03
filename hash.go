package bloomzinho

import "hash/maphash"

// masterSeed stores the seed we will use for all hashes
// this is a global instead of being unique for each filter
// to allow for consistent hashing across filters
var masterSeed = maphash.MakeSeed()

func hashS(s string) uint64 {
	return maphash.String(masterSeed, s)
}

func hashB(b []byte) uint64 {
	return maphash.Bytes(masterSeed, b)
}
