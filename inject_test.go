package inject

import (
	"errors"
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

func TestNew(t *testing.T) {
	ic := New()

	assert.IsType(t, &container.Container{}, ic)
	assert.Implements(t, new(Container), ic)
}

func TestGet(t *testing.T) {
	t.Run("get instance of specific type by name", func(t *testing.T) {
		defer Flush()

		depName := "userTest1"
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

		depName := "userTest2"
		dep := dependency.New(newUserError, name, age)

		if err := Provide(depName, dep); err != nil {
			t.Error(err)
			return
		}

		ui, err := Get[*user](depName)
		expErr := errors.New("inject: error building dependency instance: inject: error constructing `func(string, int) (*inject.user, error)`: some error")

		assert.Error(t, err)
		assert.Empty(t, ui)
		assert.Equal(t, expErr, err)
	})

	t.Run("cast type error", func(t *testing.T) {
		defer Flush()

		depName := "userTest3"
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

	depName := "userTest"
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
