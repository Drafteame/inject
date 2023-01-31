package dependency

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Drafteame/inject/dependency/mocks"
	"github.com/Drafteame/inject/types"
)

// nolint
type db interface {
	client()
}

// nolint
type database struct {
	dbname string
}

// nolint
func newDatabase(dbname string) *database {
	return &database{dbname: dbname}
}

// nolint
func (db *database) client() {}

type namer interface {
	getName() string
}

// nolint
type ager interface {
	getAge() int
}

// nolint
type userer interface {
	namer
	ager
}

// nolint
type user struct {
	name string
	age  int
	conn db
}

// nolint
func newUser(name string, age int) *user {
	return &user{name: name, age: age}
}

// nolint
func newUserConn(conn db) *user {
	return &user{conn: conn}
}

// nolint
func (u *user) getName() string {
	return u.name
}

// nolint
func (u *user) getAge() int {
	return u.age
}

func TestNewShared(t *testing.T) {
	dep := NewSingleton(func() {})

	assert.IsType(t, Dependency{}, dep)
	assert.True(t, dep.Singleton)
}

func TestDependency_IsShared(t *testing.T) {
	dep := New(func() {})
	s := NewSingleton(func() {})

	assert.True(t, s.IsSingleton())
	assert.False(t, dep.IsSingleton())
}

func TestDependency_String(t *testing.T) {
	constructor := func(string, int) {}
	dep := New(constructor, "some", 0)

	expStr := "dependency.Dependency{Factory: func(string, int), Args: [some 0]}"

	assert.Equal(t, expStr, dep.String())
}

func TestDependency_Build(t *testing.T) {
	t.Run("no arguments and no return value", func(t *testing.T) {
		constructor := func() {}

		dep := New(constructor)

		res, err := dep.Build()

		assert.NoError(t, err)
		assert.Nil(t, res)
	})

	t.Run("nil constructor", func(t *testing.T) {
		dep := New(nil)

		res, err := dep.Build()

		expErr := errors.New("inject: can't build an untyped nil")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, expErr, err)
	})

	t.Run("non-function constructors", func(t *testing.T) {
		dep := New(10)

		res, err := dep.Build()

		expErr := errors.New("inject: must provide constructor function, got `int`")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, expErr, err)
	})

	t.Run("nested deps error", func(t *testing.T) {
		dep := New(func(string) bool { return true },
			New(func() (string, error) { return "", errors.New("some") }),
		)

		res, err := dep.Build()

		expErr := errors.New("inject: error resolving argument 0 for constructor func(string) bool: inject: error constructing `func() (string, error)`: some")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, expErr, err)
	})

	t.Run("no arguments single return value", func(t *testing.T) {
		type value struct{}
		constructor := func() *value { return &value{} }

		dep := New(constructor)

		res, err := dep.Build()

		assert.NoError(t, err)
		assert.IsType(t, &value{}, res)
	})

	t.Run("no arguments and two return values", func(t *testing.T) {
		type value struct{}
		constructor := func() (*value, error) { return &value{}, nil }

		dep := New(constructor)

		res, err := dep.Build()

		assert.NoError(t, err)
		assert.IsType(t, &value{}, res)
	})

	t.Run("no arguments and two returns values - last element no error", func(t *testing.T) {
		type value struct{}
		constructor := func() (*value, int) { return &value{}, 0 }

		dep := New(constructor)

		res, err := dep.Build()

		expErr := errors.New("inject: error constructing `func() (*dependency.value, int)`: inject: last result argument of the constructor is not an error")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, expErr, err)
	})

	t.Run("no arguments and more than two return values", func(t *testing.T) {
		type value struct{}
		constructor := func() (*value, int, bool, error) { return &value{}, 0, false, nil }

		dep := New(constructor)

		res, err := dep.Build()

		assert.NoError(t, err)
		assert.IsType(t, &value{}, res)
	})

	t.Run("no arguments with error from constructor", func(t *testing.T) {
		constErr := errors.New("some error")

		type value struct{}
		constructor := func() (*value, error) { return nil, constErr }

		dep := New(constructor)

		res, err := dep.Build()

		expErr := errors.New("inject: error constructing `func() (*dependency.value, error)`: some error")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, expErr, err)
	})

	t.Run("with arguments and two return values", func(t *testing.T) {
		const name = "John Doe"
		const age = 21

		type value struct {
			name string
			age  int
		}

		constructor := func(name string, age int) (*value, error) {
			return &value{
				name: name,
				age:  age,
			}, nil
		}

		dep := New(constructor, name, age)

		res, err := dep.Build()

		assert.NoError(t, err)
		assert.IsType(t, &value{}, res)

		resValue := res.(*value)

		assert.Equal(t, name, resValue.name)
		assert.Equal(t, age, resValue.age)
	})

	t.Run("with arguments and two return values with error from constructor", func(t *testing.T) {
		const name = "John Doe"
		const age = 21

		constErr := errors.New("some error")

		type value struct {
			name string
			age  int
		}

		constructor := func(name string, age int) (*value, error) {
			return &value{
				name: name,
				age:  age,
			}, constErr
		}

		dep := New(constructor, name, age)

		res, err := dep.Build()

		expErr := errors.New("inject: error constructing `func(string, int) (*dependency.value, error)`: some error")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, expErr, err)
	})

	t.Run("with arguments and error by wrong argument order or type", func(t *testing.T) {
		const name = "John Doe"
		const age = 21

		constErr := errors.New("some error")

		type value struct {
			name string
			age  int
		}

		constructor := func(name string, age int) (*value, error) {
			return &value{
				name: name,
				age:  age,
			}, constErr
		}

		dep := New(constructor, age, name)

		res, err := dep.Build()

		expErr := errors.New("inject: using int as type string on constructor `func(string, int) (*dependency.value, error)`")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, expErr, err)
	})

	t.Run("with arguments and error by wrong number passed to constructor", func(t *testing.T) {
		const name = "John Doe"

		constErr := errors.New("some error")

		type value struct {
			name string
			age  int
		}

		constructor := func(name string, age int) (*value, error) {
			return &value{
				name: name,
				age:  age,
			}, constErr
		}

		dep := New(constructor, name)

		res, err := dep.Build()

		expErr := errors.New("inject: invalid argument length for constructor `func(string, int) (*dependency.value, error)`, got 1 (need 2)")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, expErr, err)
	})

	t.Run("with function as argument", func(t *testing.T) {
		const name = "John Smith"
		const age = 21
		var callback = func() int { return age }

		type value struct {
			name string
			age  int
		}

		var construct = func(name string, age int) *value {
			return &value{name: name, age: age}
		}

		dep := New(construct, name, callback)

		val, err := dep.Build()

		if assert.NoError(t, err) && assert.IsType(t, &value{}, val) {
			obj := val.(*value)

			assert.Equal(t, age, obj.age)
			assert.Equal(t, name, obj.name)
		}
	})

	t.Run("with nested dependency as argument", func(t *testing.T) {
		const name = "John Smith"
		const age = 21

		type value struct {
			userService userer
		}

		var construct = func(userService userer) *value {
			return &value{userService: userService}
		}

		dep := New(
			construct,
			New(newUser, name, age),
		)

		val, err := dep.Build()

		if assert.NoError(t, err) && assert.IsType(t, &value{}, val) {
			obj := val.(*value)

			assert.Equal(t, age, obj.userService.getAge())
			assert.Equal(t, name, obj.userService.getName())
		}
	})

	t.Run("with nil value as argument over interface", func(t *testing.T) {
		type value struct {
			userService userer
		}

		var construct = func(userService userer) *value {
			return &value{userService: userService}
		}

		dep := New(
			construct,
			nil,
		)

		val, err := dep.Build()

		if assert.NoError(t, err) && assert.IsType(t, &value{}, val) {
			obj := val.(*value)

			assert.Nil(t, obj.userService)
		}
	})

	t.Run("with nil value as argument over pointer", func(t *testing.T) {
		type value struct {
			userService userer
		}

		var construct = func(userService *user) *value {
			return &value{userService: userService}
		}

		dep := New(
			construct,
			nil,
		)

		val, err := dep.Build()

		if assert.NoError(t, err) && assert.IsType(t, &value{}, val) {
			obj := val.(*value)

			assert.Nil(t, obj.userService)
		}
	})

	t.Run("with injectable dependency", func(t *testing.T) {
		injectDepName := types.Symbol("inject")
		injectDep := Inject(injectDepName)
		injectDepValue := "some"

		ic := mocks.NewContainer(t)
		ic.On("Get", injectDepName).Return(injectDepValue, nil)

		injectedValue := ""

		dep := New(func(name string) int {
			injectedValue = name
			return 1
		}, injectDep)

		res, err := dep.SetContainer(ic).Build()

		assert.NoError(t, err)
		assert.Equal(t, (any)(1), res)
		assert.Equal(t, injectDepValue, injectedValue)
	})
}
