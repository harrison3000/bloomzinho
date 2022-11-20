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
