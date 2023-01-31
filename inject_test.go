package inject

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Drafteame/inject/container"
	"github.com/Drafteame/inject/dependency"
	"github.com/Drafteame/inject/types"
)

const name = "John"

const age = 21

type user struct {
	name string
	age  int
	db   *sql.DB
}

func newUser(name string, age int) *user {
	return &user{
		age:  age,
		name: name,
	}
}

func newUserError(_ string, _ int) (*user, error) {
	return nil, errors.New("some error")
}

func newUserWithDB(db *sql.DB) *user {
	return &user{db: db}
}

func newDB() *sql.DB {
	return &sql.DB{}
}

func TestNew(t *testing.T) {
	ic := New()

	assert.IsType(t, &container.Container{}, ic)
	assert.Implements(t, new(Container), ic)
}

func TestGet(t *testing.T) {
	t.Run("get instance of specific type by name", func(t *testing.T) {
		defer Flush()

		depName := types.Symbol("userTest1")
		dep := dependency.New(newUser, name, age)

		if err := Provide(depName, dep); err != nil {
			t.Error(err)
			return
		}

		ui, err := Get[*user](depName)

		assert.NoError(t, err)
		assert.NotEmpty(t, ui)
		assert.Equal(t, ui.age, age)
		assert.Equal(t, ui.name, name)
	})

	t.Run("get instance from container error", func(t *testing.T) {
		defer Flush()

		depName := types.Symbol("userTest2")
		dep := dependency.New(newUserError, name, age)

		if err := Provide(depName, dep); err != nil {
			t.Error(err)
			return
		}

		ui, err := Get[*user](string(depName))
		expErr := errors.New("inject: error building dependency instance: inject: error constructing `func(string, int) (*inject.user, error)`: some error")

		assert.Error(t, err)
		assert.Empty(t, ui)
		assert.Equal(t, expErr, err)
	})

	t.Run("cast type error", func(t *testing.T) {
		defer Flush()

		depName := types.Symbol("userTest3")
		dep := dependency.New(newUser, name, age)

		if err := Provide(depName, dep); err != nil {
			t.Error(err)
			return
		}

		ui, err := Get[string](depName)
		expErr := errors.New("inject: error casting instance of `userTest3` dependency to `string`")

		assert.Error(t, err)
		assert.Empty(t, ui)
		assert.Equal(t, expErr, err)
	})
}

func TestInvoke(t *testing.T) {
	defer Flush()

	depName := types.Symbol("userTest")
	dep := dependency.New(newUser, name, age)

	if err := Provide(depName, dep); err != nil {
		t.Error(err)
		return
	}

	type args struct {
		types.In
		User *user `inject:"name=userTest"`
	}

	called := false

	invoker := func(in args) {
		if assert.NotNil(t, in.User) {
			assert.Equal(t, in.User.age, age)
			assert.Equal(t, in.User.name, name)
		}

		called = true
	}

	err := Invoke(invoker)

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestSingleton(t *testing.T) {
	t.Run("should register a raw factory singleton instance", func(t *testing.T) {
		defer Flush()

		factoryName := "test"

		err := Singleton(factoryName, newUser, name, age)

		assert.NoError(t, err)
	})

	t.Run("should register a singleton from singleton dependency", func(t *testing.T) {
		defer Flush()

		factoryName := "test"

		dep := dependency.NewSingleton(newUser, name, age)

		err := Singleton(factoryName, dep)

		assert.NoError(t, err)
	})

	t.Run("should register a singleton from dependency", func(t *testing.T) {
		defer Flush()

		factoryName := "test"

		dep := dependency.New(newUser, name, age)

		err := Singleton(factoryName, dep)

		assert.NoError(t, err)
	})

	t.Run("should register a singleton from raw function and nested dependencies", func(t *testing.T) {
		defer Flush()

		depName := "db"

		if err := Singleton(depName, newDB); err != nil {
			t.Error(err)
			return
		}

		factoryName := "test"

		if err := Singleton(factoryName, newUserWithDB, Dep(depName)); err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("error when no dependency.Depdndendency instance or raw function is registered", func(t *testing.T) {
		defer Flush()

		err := Singleton("name", dependency.Injectable{})

		expErr := fmt.Errorf("factory parameter should be a function or a dependency.Dependency instance")

		if assert.Error(t, err) {
			assert.Equal(t, expErr, err)
		}
	})
}

func TestProvide(t *testing.T) {
	t.Run("should register a raw factory instance", func(t *testing.T) {
		defer Flush()

		factoryName := "test"

		err := Provide(factoryName, newUser, name, age)

		assert.NoError(t, err)
	})

	t.Run("should register a singleton from singleton dependency", func(t *testing.T) {
		defer Flush()

		factoryName := "test"

		dep := dependency.NewSingleton(newUser, name, age)

		err := Provide(factoryName, dep)

		assert.NoError(t, err)
	})

	t.Run("should register a singleton from dependency", func(t *testing.T) {
		defer Flush()

		factoryName := "test"

		dep := dependency.New(newUser, name, age)

		err := Provide(factoryName, dep)

		assert.NoError(t, err)
	})

	t.Run("should register a singleton from raw function and nested dependencies", func(t *testing.T) {
		defer Flush()

		depName := "db"

		if err := Provide(depName, newDB); err != nil {
			t.Error(err)
			return
		}

		factoryName := "test"

		if err := Provide(factoryName, newUserWithDB, Dep(depName)); err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("error when no dependency.Depdndendency instance or raw function is registered", func(t *testing.T) {
		defer Flush()

		err := Provide("name", dependency.Injectable{})

		expErr := fmt.Errorf("factory parameter should be a function or a dependency.Dependency instance")

		if assert.Error(t, err) {
			assert.Equal(t, expErr, err)
		}
	})
}
