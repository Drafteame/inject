package container

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Drafteame/inject/dependency"
	"github.com/Drafteame/inject/types"
)

func TestContainer_Provide(t *testing.T) {
	t.Run("provide simple dependency", func(t *testing.T) {
		const name = "John Smith"
		const age = 21

		ic := New()

		userDepName := types.Symbol("userDep")
		userDep := dependency.New(newUser, name, age)

		err := ic.Provide(userDepName, userDep)

		if assert.NoError(t, err) {
			assert.Len(t, ic.deps, 1)

			dep := ic.deps[userDepName]

			assert.Equal(t, userDep.String(), dep.String())
		}
	})

	t.Run("provide duplicated dependency", func(t *testing.T) {
		const name = "John Smith"
		const age = 21

		ic := New()

		userDepName := types.Symbol("usersService")
		userDep := dependency.New(newUser, name, age)

		if err := ic.Provide(userDepName, userDep); err != nil {
			t.Error(err)
			return
		}

		err := ic.Provide(userDepName, userDep)

		expErr := fmt.Errorf("inject: duplicated dependency name `%s`", userDepName)

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("provide dependency with no return value constructor", func(t *testing.T) {
		ic := New()

		depName := types.Symbol("test")
		dep := dependency.New(func() {})

		err := ic.Provide(depName, dep)
		expErr := fmt.Errorf("inject: dependency factory should return at least one return type: dependency.Dependency{Factory: func(), Args: []}")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("provide singleton dependency", func(t *testing.T) {
		ic := New()

		depName := types.Symbol("test")
		dep := dependency.NewSingleton(newDriver, "test")

		err := ic.Provide(depName, dep)

		assert.NoError(t, err)

		assert.True(t, ic.deps[depName].IsSingleton())
	})

	t.Run("provide duplicated shared named dependency", func(t *testing.T) {
		const name = "John Smith"
		const age = 21

		ic := New()

		userDepName := types.Symbol("usersService")
		userDep := dependency.NewSingleton(newUser, name, age)
		if err := ic.Provide(userDepName, userDep); err != nil {
			t.Error(err)
			return
		}

		err := ic.Provide(userDepName, userDep)

		expErr := fmt.Errorf("inject: duplicated dependency name `%s`", userDepName)

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("provide singleton dependency with no return value constructor", func(t *testing.T) {
		ic := New()

		depName := types.Symbol("test")
		dep := dependency.NewSingleton(func() {})

		err := ic.Provide(depName, dep)
		expErr := fmt.Errorf("inject: dependency factory should return at least one return type: dependency.Dependency{Factory: func(), Args: []}")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("provide dependency with empty name", func(t *testing.T) {
		const name = "John Smith"
		const age = 21

		ic := New()

		userDepName := types.Symbol("")
		userDep := dependency.New(newUser, name, age)

		err := ic.Provide(userDepName, userDep)

		expErr := errors.New("inject: dependency name cannot be empty")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("wrong container initialization solved", func(t *testing.T) {
		const name = "John Smith"
		const age = 21

		ic := &Container{}

		userDepName := types.Symbol("some")
		userDep := dependency.New(newUser, name, age)

		err := ic.Provide(userDepName, userDep)

		assert.NoError(t, err)
		assert.NotEmpty(t, ic.deps[userDepName])
	})
}
