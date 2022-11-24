package bloomzinho

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntersectSimple(t *testing.T) {
	a := NewFilter(32768, 2)
	b := NewFilter(32768, 2)

	a.AddString("Vem meu amor")
	b.AddString("Me tirar da solidão")

	assert.False(t, a.Intersects(b))

	a.AddString("Hello")
	b.AddString("Hello")

	assert.True(t, a.Intersects(b))
}

func BenchmarkIntersects(b *testing.B) {
	a := NewFilter(32768, 2)
	f := NewFilter(32768, 2)

	a.AddString("Ehh")
	//f.AddString("Ehh")
	f.AddString("Ohh")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		a.Intersects(f)
	}
}
