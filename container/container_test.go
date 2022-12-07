package container

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Drafteame/inject/dependency"
)

func TestContainer_Flush(t *testing.T) {
	depName := "test"

	ic := New()

	if err := ic.Provide(depName, dependency.New(func() int { return 10 })); err != nil {
		t.Error(err)
		return
	}

	ic.Flush()

	assert.Empty(t, ic.deps)
	assert.Empty(t, ic.solvedDeps)
}
