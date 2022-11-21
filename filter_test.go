package bloomzinho

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrivial(t *testing.T) {
	f := NewFilter(256, 3)

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

func BenchmarkTrivial(b *testing.B) {
	f := NewFilter(4096, 5)

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
		m := make([]byte, 12_000)
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
		f := NewFilter(256, 2)
		for i := 0; i < 16; i++ {
			st := gimme(16)
			f.AddString(st)
		}
		b.StartTimer()

		for i := 0; i < 1000; i++ {
			st := gimme(10)
			if f.LookupString(st) {
				positive++
			}
			total++
		}
	}

	b.ReportMetric((positive/total)*100, "%FalsePositives")
}
