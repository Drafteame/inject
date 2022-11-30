package inject

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Drafteame/inject/dependency"
)

func TestContainer_Provide(t *testing.T) {
	t.Run("provide simple dependency", func(t *testing.T) {
		const name = "John Smith"
		const age = 21

		inject := &container{}

		userDep := dependency.New(newUser, name, age)

		err := inject.Provide(userDep)

		if assert.NoError(t, err) {
			assert.Len(t, inject.deps, 1)

			dep := inject.deps[reflect.TypeOf(&user{})]

			assert.Equal(t, userDep.String(), dep.String())
		}
	})

	t.Run("provide dependency with As option", func(t *testing.T) {
		const name = "John Smith"
		const age = 21

		inject := &container{}

		userDep := dependency.New(newUser, name, age)

		err := inject.Provide(userDep, As(new(namer)))

		if assert.NoError(t, err) {
			assert.Len(t, inject.deps, 2)

			dep := inject.deps[reflect.TypeOf(&user{})]
			assert.Equal(t, userDep.String(), dep.String())

			depAlias := inject.deps[reflect.TypeOf(new(namer)).Elem()]
			assert.Equal(t, userDep.String(), depAlias.String())
		}
	})

	t.Run("provide dependency with Name option", func(t *testing.T) {
		const name = "John Smith"
		const age = 21
		const depName = "usersService"

		inject := &container{}

		userDep := dependency.New(newUser, name, age)

		err := inject.Provide(userDep, Name(depName))

		if assert.NoError(t, err) {
			assert.Len(t, inject.deps, 1)
			assert.Len(t, inject.depsByName, 1)

			dep := inject.deps[reflect.TypeOf(&user{})]
			assert.Equal(t, userDep.String(), dep.String())

			depAlias := inject.depsByName[depName]
			assert.Equal(t, userDep.String(), depAlias.String())
		}
	})

	t.Run("provide duplicated named dependency", func(t *testing.T) {
		const name = "John Smith"
		const age = 21
		const depName = "usersService"

		inject := &container{}

		userDep := dependency.New(newUser, name, age)

		if err := inject.Provide(userDep, Name(depName)); err != nil {
			t.Error(err)
			return
		}

		err := inject.Provide(userDep, Name(depName))

		expErr := fmt.Errorf("inject: duplicated dependency name `%s`", depName)

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("provide dependency with wrong As alias option", func(t *testing.T) {
		const name = "John Smith"
		const age = 21

		inject := &container{}

		userDep := dependency.New(newUser, name, age)

		err := inject.Provide(userDep, As(user{}))

		expErr := fmt.Errorf("inject: alias option value should be a pointer to an interface")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("provide dependency with no return value constructor", func(t *testing.T) {
		construct := func() {}
		inject := &container{}

		err := inject.Provide(dependency.New(construct))
		expErr := fmt.Errorf("inject: can't provide a dependency constructor with no return types: dependency.Dependency{Constructor: func(), Args: []}")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("provide shared dependency", func(t *testing.T) {
		const dbname = "test"
		inject := newContainer()

		err := inject.Provide(dependency.NewShared(newDriver, dbname))

		assert.NoError(t, err)

		assert.Len(t, inject.shared.deps, 1)
	})

	t.Run("provide shared dependency with alias", func(t *testing.T) {
		const dbname = "test"
		inject := newContainer()

		dep := dependency.NewShared(newDriver, dbname)

		err := inject.Provide(dep, As(new(database)))

		assert.NoError(t, err)

		assert.Len(t, inject.shared.deps, 2)
	})

	t.Run("provide shared dependency with name", func(t *testing.T) {
		const dbname = "test"
		const depName = "driver"

		inject := newContainer()

		dep := dependency.NewShared(newDriver, dbname)

		err := inject.Provide(dep, Name(depName))

		assert.NoError(t, err)

		assert.Len(t, inject.shared.deps, 1)
		assert.Len(t, inject.shared.depsByName, 1)
	})

	t.Run("provide duplicated shared named dependency", func(t *testing.T) {
		const name = "John Smith"
		const age = 21
		const depName = "usersService"

		inject := newContainer()

		userDep := dependency.NewShared(newUser, name, age)

		if err := inject.Provide(userDep, Name(depName)); err != nil {
			t.Error(err)
			return
		}

		err := inject.Provide(userDep, Name(depName))

		expErr := fmt.Errorf("inject: duplicated dependency name `%s`", depName)

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("provide shared dependency with wrong As alias option", func(t *testing.T) {
		const name = "John Smith"
		const age = 21

		inject := newContainer()

		userDep := dependency.NewShared(newUser, name, age)

		err := inject.Provide(userDep, As(user{}))

		expErr := fmt.Errorf("inject: alias option value should be a pointer to an interface")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("provide shared dependency with no return value constructor", func(t *testing.T) {
		construct := func() {}
		inject := newContainer()

		err := inject.Provide(dependency.NewShared(construct))
		expErr := fmt.Errorf("inject: can't provide a dependency constructor with no return types: dependency.Dependency{Constructor: func(), Args: []}")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})
}
