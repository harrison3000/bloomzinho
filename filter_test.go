package bloomzinho

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

// TODO refactor all the tests

func TestTrivial(t *testing.T) {
	f, e := NewFilter(256, 3)
	assert.NoError(t, e)

	f.AddString("hellow")
	f.AddString("okay")
	f.AddBytes([]byte("bye"))

	//TODO sometimes there are false positives and the test fails...
	//needs to make this test probabilistic somehow
	assert.False(t, f.LookupString("eita"))
	assert.False(t, f.LookupString("parafuso"))
	assert.False(t, f.LookupBytes([]byte("água de coco")))

	assert.True(t, f.LookupString("bye"))
	assert.True(t, f.LookupString("hellow"))
	assert.True(t, f.LookupBytes([]byte("okay")))
}

func TestSingleBit(t *testing.T) {
	//this test ensures that the loop is not off by one
	f, _ := NewFilter(256, 1)

	f.AddString("hi")

	assert.False(t, f.LookupString("hello"))
	assert.True(t, f.LookupString("hi"))
}

func TestLookup(t *testing.T) {
	f, _ := NewFilter(257, 1)

	f.set(63)
	f.set(94)
	f.set(260) //we give some bits more, the number is always a multiple of 64

	assert.False(t, f.lookup(25))
	assert.False(t, f.lookup(0))
	assert.False(t, f.lookup(64))
	assert.True(t, f.lookup(63))
	assert.True(t, f.lookup(94))
	assert.True(t, f.lookup(260))
}

func TestBigHashShuffle(t *testing.T) {
	//this test ensures the shuffle really does something

	//4 groups of 16bit numbers, but on different order
	a := uint64(0x1234_5678_9abc_def0)
	b := uint64(0x5678_def0_9abc_1234)

	f, _ := NewFilter(1<<16, 4)

	hToIdx := func(hash uint64) (idx []uint) {
		f.hashTransform(hash, func(u uint) {
			idx = append(idx, u)
		})
		slices.Sort(idx)
		return
	}

	a4 := hToIdx(a)
	b4 := hToIdx(b)

	//even though the hashes are different,
	//the indexes end up being the same thing...
	assert.Equal(t, a4, b4)

	f, _ = NewFilter(1<<16, 8)
	a8 := hToIdx(a)
	b8 := hToIdx(b)

	//..however the indexes after the shuffle must be different
	assert.NotEqual(t, a8, b8, "the bit shuffle is not working")
}

func TestError(t *testing.T) {
	f, e := NewFilter(0, 0)

	assert.Nil(t, f)
	assert.Error(t, e)
}

func BenchmarkTrivial(b *testing.B) {
	f, _ := NewFilter(9585059, 7)

	f.AddString("whatever")

	for i := 0; i < b.N; i++ {
		f.LookupString("I eat a log")
		f.LookupString("we be blue")
		f.LookupString("I think it was a bee")
		f.LookupString("and I flee a salami")
	}
}

func BenchmarkFalsePositive(b *testing.B) {
	var alea string
	si := 0
	seed := func() {
		m := make([]byte, 100_000)
		rand.Read(m)
		alea = string(m)
		si = 0
	}
	gimme := func(n int) string {
		st := si
		si = st + n
		return alea[st:si]
	}

	b.ResetTimer()

	positive := 0.0
	total := 0.0
	for t := 0; t < b.N; t++ {
		b.StopTimer()
		seed()
		f, _ := NewFilter(4094, 7)
		for i := 0; i < 128; i++ {
			st := gimme(16)
			f.AddString(st)
		}
		b.StartTimer()

		for i := 0; i < 10000; i++ {
			st := gimme(1 + i%8)
			if f.LookupString(st) {
				positive++
			}
			total++
		}
	}

	b.ReportMetric((positive/total)*100, "%FalsePositives")
}

func TestNoAlloc(t *testing.T) {
	//this test ensures that we don't make any alocations
	//other than when making the filter, of course

	f, _ := NewFilter(256, 1)
	fb, _ := NewFilter(256, 1)

	a := testing.AllocsPerRun(5, func() {
		f.AddBytes([]byte("high"))
		fb.AddString("low")

		f.LookupBytes([]byte("why?"))
		fb.LookupString("I don't know")

		f.Intersects(fb)
	})

	assert.Zero(t, a, "Allocations made on methods that should be zero alloc")
}
