package dependency

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Drafteame/inject/dependency/mocks"
	"github.com/Drafteame/inject/types"
)

func TestInject(t *testing.T) {
	name := types.Symbol("test")
	i := Inject(name)

	assert.IsType(t, Injectable{}, i)
	assert.Equal(t, name, i.name)
}

func TestInjectable_IsSingleton(t *testing.T) {
	name := types.Symbol("test")
	i := Inject(name)

	assert.False(t, i.IsSingleton())
}

func TestInjectable_SetContainer(t *testing.T) {
	ic := mocks.NewContainer(t)

	name := types.Symbol("test")
	i := Inject(name).SetContainer(ic).(Injectable)

	assert.NotNil(t, i.container)
	assert.Same(t, ic, i.container)
}

func TestInjectable_Build(t *testing.T) {
	t.Run("resolve build from container", func(t *testing.T) {
		depName := types.Symbol("test")
		dep := Inject(depName)

		ic := mocks.NewContainer(t)
		ic.On("Get", depName).Return("some", nil)

		res, err := dep.SetContainer(ic).Build()

		assert.NoError(t, err)
		assert.Equal(t, (any)("some"), res)
	})

	t.Run("error by empty container", func(t *testing.T) {
		depName := types.Symbol("test")
		dep := Inject(depName)

		_, err := dep.Build()
		expErr := errors.New("inject: [internal-error] no container provided")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})
}
