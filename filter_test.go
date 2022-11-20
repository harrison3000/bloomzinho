package bloomzinho

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrivial(t *testing.T) {
	f := NewFilter(256, 3)

	f.AddString("hellow")
	f.AddString("okay")
	f.AddString("bye")

	//TODO sometimes there are false positives and the test fails...
	//needs to make this test probabilistic somehow
	assert.False(t, f.LookupString("eita"))
	assert.False(t, f.LookupString("parafuso"))

	assert.True(t, f.LookupString("bye"))
	assert.True(t, f.LookupString("hellow"))
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
