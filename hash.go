package bloomzinho

import "hash/maphash"

var masterSeed = maphash.MakeSeed()

func hashS(s string) uint64 {
	return maphash.String(masterSeed, s)
}

func hashB(b []byte) uint64 {
	return maphash.Bytes(masterSeed, b)
}
