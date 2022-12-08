package bloomzinho

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntersectSimple(t *testing.T) {
	a, _ := NewFilter(32768, 1)
	b, _ := NewFilter(32768, 1)

	a.AddString("Vem meu amor")
	b.AddString("Me tirar da solidão")

	assert.False(t, a.Intersects(b))

	a.AddString("Hello")
	b.AddString("Hello")

	assert.True(t, a.Intersects(b))
}

func BenchmarkIntersects(b *testing.B) {
	a, _ := NewFilter(32768, 1)
	f, _ := NewFilter(32768, 1)

	a.AddString("Ehh")
	//f.AddString("Ehh")
	f.AddString("Ohh")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		a.Intersects(f)
	}
}

//TODO examples
//TODO newunion and newintersection tests
