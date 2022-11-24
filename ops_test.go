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
