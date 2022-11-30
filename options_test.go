package inject

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAs(t *testing.T) {
	t.Run("alias from interface", func(t *testing.T) {
		opts := injectOptions{}

		as := As(new(namer))

		err := as(&opts)

		assert.NoError(t, err)
		assert.Len(t, opts.aliases, 1)
	})

	t.Run("alias from pointer", func(t *testing.T) {
		opts := injectOptions{}

		as := As(&user{})

		err := as(&opts)
		expErr := fmt.Errorf("inject: alias option value should be a pointer to an interface")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})
}

func TestName(t *testing.T) {
	t.Run("set provider name", func(t *testing.T) {
		opts := injectOptions{}

		name := Name("test")

		err := name(&opts)

		assert.NoError(t, err)
		assert.Equal(t, "test", *opts.name)
	})

	t.Run("error empty name", func(t *testing.T) {
		opts := injectOptions{}

		name := Name("")

		err := name(&opts)

		expErr := fmt.Errorf("inject: name cannot be empty if provided")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})
}
