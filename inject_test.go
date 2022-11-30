package inject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	inject := Get()

	assert.NotNil(t, inject)
	assert.IsType(t, &container{}, inject)
	assert.Implements(t, new(Injector), inject)
}

func TestNew(t *testing.T) {
	inject := Get()
	inject2 := New()

	assert.NotNil(t, inject)
	assert.IsType(t, &container{}, inject)
	assert.Implements(t, new(Injector), inject)

	assert.NotNil(t, inject2)
	assert.IsType(t, &container{}, inject2)
	assert.Implements(t, new(Injector), inject2)

	assert.NotSame(t, inject, inject2)
}
