package inject

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Drafteame/inject/dependency"
)

func TestContainer_Invoke(t *testing.T) {
	t.Run("invoke no dependency function", func(t *testing.T) {
		inject := newContainer()
		called := false

		invoker := func() {
			called = true
		}

		err := inject.Invoke(invoker)

		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("invoke with dependencies", func(t *testing.T) {
		const name = "John Smith"
		const age = 21

		inject := newContainer()
		called := false

		type args struct {
			In
			UserService *user
		}

		invoker := func(in args) {
			if assert.NotNil(t, in.UserService) {
				assert.Equal(t, name, in.UserService.getName())
				assert.Equal(t, age, in.UserService.getAge())
			}

			called = true
		}

		userDep := dependency.New(newUser, name, age)

		if err := inject.Provide(userDep); err != nil {
			t.Error(err)
			return
		}

		err := inject.Invoke(invoker)

		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("invoke with dependencies and dependency alias", func(t *testing.T) {
		const name = "John Smith"
		const age = 21

		inject := newContainer()
		called := false

		type args struct {
			In
			UserService userer
		}

		invoker := func(in args) {
			if assert.NotNil(t, in.UserService) {
				assert.Equal(t, name, in.UserService.getName())
				assert.Equal(t, age, in.UserService.getAge())
			}

			called = true
		}

		userDep := dependency.New(newUser, name, age)

		if err := inject.Provide(userDep, As(new(userer))); err != nil {
			t.Error(err)
			return
		}

		err := inject.Invoke(invoker)

		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("invoke with dependency and dependency name", func(t *testing.T) {
		const name = "John Smith"
		const age = 21
		const depName = "usersService"

		inject := newContainer()
		called := false

		type args struct {
			In
			UserService userer `inject:"name=usersService"`
		}

		invoker := func(in args) {
			if assert.NotNil(t, in.UserService) {
				assert.Equal(t, name, in.UserService.getName())
				assert.Equal(t, age, in.UserService.getAge())
			}

			called = true
		}

		userDep := dependency.New(newUser, name, age)

		if err := inject.Provide(userDep, Name(depName)); err != nil {
			t.Error(err)
			return
		}

		err := inject.Invoke(invoker)

		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("invoke with dependency and dependency name and dependency alias", func(t *testing.T) {
		const name = "John Smith"
		const age = 21
		const depName = "usersService"

		inject := newContainer()
		called := false

		type args struct {
			In
			UserService      userer `inject:"name=usersService"`
			UserServiceAlias userer
		}

		invoker := func(in args) {
			if assert.NotNil(t, in.UserService) {
				assert.Equal(t, name, in.UserService.getName())
				assert.Equal(t, age, in.UserService.getAge())

				assert.Equal(t, name, in.UserServiceAlias.getName())
				assert.Equal(t, age, in.UserServiceAlias.getAge())
			}

			called = true
		}

		userDep := dependency.New(newUser, name, age)

		if err := inject.Provide(userDep, Name(depName), As(new(userer))); err != nil {
			t.Error(err)
			return
		}

		err := inject.Invoke(invoker)

		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("invoke with dependency and optional injection", func(t *testing.T) {
		const name = "John Smith"
		const age = 21
		const depName = "usersService"

		inject := newContainer()
		called := false

		type args struct {
			In
			UserService namer `inject:"optional"`
		}

		invoker := func(in args) {
			assert.Nil(t, in.UserService)

			called = true
		}

		userDep := dependency.New(newUser, name, age)

		if err := inject.Provide(userDep, Name(depName)); err != nil {
			t.Error(err)
			return
		}

		err := inject.Invoke(invoker)

		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("invoke with empty inject tag dependency", func(t *testing.T) {
		const name = "John Smith"
		const age = 21

		inject := newContainer()
		called := false

		type args struct {
			In
			UserService userer `inject:""`
		}

		invoker := func(in args) {
			if assert.NotNil(t, in.UserService) {
				assert.Equal(t, name, in.UserService.getName())
				assert.Equal(t, age, in.UserService.getAge())
			}

			called = true
		}

		userDep := dependency.New(newUser, name, age)

		if err := inject.Provide(userDep, As(new(userer))); err != nil {
			t.Error(err)
			return
		}

		err := inject.Invoke(invoker)

		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("invoke with nil invoker", func(t *testing.T) {
		inject := newContainer()

		err := inject.Invoke(nil)

		expErr := fmt.Errorf("inject: can't invoke nil constructor")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("invoke with non-function invoker", func(t *testing.T) {
		inject := newContainer()

		err := inject.Invoke(10)

		expErr := fmt.Errorf("inject: can't invoke a non-function constructor")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("invoke with non In embedded struct", func(t *testing.T) {
		inject := newContainer()

		type args struct{}

		invoker := func(in args) {}

		err := inject.Invoke(invoker)

		expErr := fmt.Errorf("inject: struct doesn't embed `inject.In` struct")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("invoke with invoker returning non error value", func(t *testing.T) {
		inject := newContainer()

		type args struct {
			In
		}

		called := false

		invoker := func(in args) bool {
			called = true
			return called
		}

		err := inject.Invoke(invoker)

		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("invoke with invoker returning nil error value", func(t *testing.T) {
		inject := newContainer()

		type args struct {
			In
		}

		called := false

		invoker := func(in args) error {
			called = true
			return nil
		}

		err := inject.Invoke(invoker)

		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("invoke with invoker returning nil error value", func(t *testing.T) {
		inject := newContainer()

		type args struct {
			In
		}

		called := false
		expErr := fmt.Errorf("invoke with invoker returning nil error value")

		invoker := func(in args) error {
			called = true
			return expErr
		}

		err := inject.Invoke(invoker)

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
		assert.True(t, called)
	})

	t.Run("invoke with error providing unnamed dependency", func(t *testing.T) {
		inject := newContainer()

		type some interface {
			someMethod()
		}

		type args struct {
			In
			UserService some
		}

		invoker := func(in args) {}

		err := inject.Invoke(invoker)

		expErr := fmt.Errorf("inject: can't provide dependency of type `inject.some` to In receiver")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("invoke with error providing named dependency", func(t *testing.T) {
		inject := newContainer()

		type args struct {
			In
			UserService *user `inject:"name=usersService"`
		}

		invoker := func(in args) {}

		err := inject.Invoke(invoker)

		expErr := fmt.Errorf("inject: can't provide dependency of type `*inject.user` and name usersService to In receiver")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("invoke error resolving dependency three", func(t *testing.T) {
		inject := newContainer()

		type args struct {
			In
			UserService *user
		}

		invoker := func(in args) {}

		dep := dependency.New(func() (*user, error) { return nil, errors.New("some") })

		if err := inject.Provide(dep); err != nil {
			t.Error(err)
			return
		}

		err := inject.Invoke(invoker)

		expErr := fmt.Errorf("inject: error constructing `func() (*inject.user, error)`: some")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("invoke with shared dependency from type", func(t *testing.T) {
		inject := newContainer()

		const driverName = "some"

		sharedDep := dependency.NewShared(newDriver, driverName)
		alias := new(database)

		if err := inject.Provide(sharedDep, As(alias)); err != nil {
			t.Error(err)
			return
		}

		usersRepo := dependency.New(newUserWithDriver, dependency.FromSharedWithType(alias))
		if err := inject.Provide(usersRepo); err != nil {
			t.Error(err)
			return
		}

		type args struct {
			In
			UserRepo *user
		}

		called := false

		invoker := func(in args) {
			assert.Equal(t, driverName, in.UserRepo.getDb().client())
			called = true
		}

		err := inject.Invoke(invoker)

		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("invoke with shared dependency from type in multiple targets", func(t *testing.T) {
		inject := newContainer()

		const driverName = "some"

		sharedDep := dependency.NewShared(newDriver, driverName)

		if err := inject.Provide(sharedDep, As(new(database))); err != nil {
			t.Error(err)
			return
		}

		usersRepo := dependency.New(newUserWithDriver, dependency.FromSharedWithType(new(database)))
		if err := inject.Provide(usersRepo); err != nil {
			t.Error(err)
			return
		}

		todoRepo := dependency.New(newTodo, dependency.FromSharedWithType(new(database)))
		if err := inject.Provide(todoRepo); err != nil {
			t.Error(err)
			return
		}

		type args struct {
			In
			UserRepo *user
			TodoRepo *todo
		}

		called := false

		invoker := func(in args) {
			assert.Equal(t, driverName, in.UserRepo.getDb().client())
			assert.Equal(t, driverName, in.TodoRepo.db.client())
			assert.Same(t, in.TodoRepo.db, in.UserRepo.getDb())
			called = true
		}

		err := inject.Invoke(invoker)

		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("invoke with shared dependency from name", func(t *testing.T) {
		inject := newContainer()

		const dbName = "main"
		const driverName = "some"

		sharedDep := dependency.NewShared(newDriver, dbName)

		if err := inject.Provide(sharedDep, Name(driverName)); err != nil {
			t.Error(err)
			return
		}

		usersRepo := dependency.New(newUserWithDriver, dependency.FromSharedWithName(driverName))
		if err := inject.Provide(usersRepo); err != nil {
			t.Error(err)
			return
		}

		type args struct {
			In
			UserRepo *user
		}

		called := false

		invoker := func(in args) {
			assert.Equal(t, dbName, in.UserRepo.getDb().client())
			called = true
		}

		err := inject.Invoke(invoker)

		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("invoke with unregistered shared dependency by type", func(t *testing.T) {
		inject := newContainer()

		const dbName = "main"

		sharedDep := dependency.NewShared(newDriver, dbName)

		if err := inject.Provide(sharedDep); err != nil {
			t.Error(err)
			return
		}

		usersRepo := dependency.New(newUserWithDriver, dependency.FromSharedWithType(new(database)))
		if err := inject.Provide(usersRepo); err != nil {
			t.Error(err)
			return
		}

		type args struct {
			In
			UserRepo *user
		}

		invoker := func(in args) {}

		err := inject.Invoke(invoker)
		expErr := errors.New("inject: error resolving argument 0 for constructor func(inject.database) *inject.user: inject: can't solve shared dependency by type inject.database")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("invoke with shared dependency by type that return error", func(t *testing.T) {
		inject := newContainer()

		const dbName = "main"

		errBuild := errors.New("some")

		sharedDep := dependency.NewShared(func(dbname string) (*driver, error) { return nil, errBuild }, dbName)

		if err := inject.Provide(sharedDep, As(new(database))); err != nil {
			t.Error(err)
			return
		}

		usersRepo := dependency.New(newUserWithDriver, dependency.FromSharedWithType(new(database)))
		if err := inject.Provide(usersRepo); err != nil {
			t.Error(err)
			return
		}

		type args struct {
			In
			UserRepo *user
		}

		invoker := func(in args) {}

		err := inject.Invoke(invoker)
		expErr := errors.New("inject: error resolving argument 0 for constructor func(inject.database) *inject.user: inject: can't solve shared dependency by type inject.database: inject: error constructing `func(string) (*inject.driver, error)`: some")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("invoke with shared dependency on multiple targets", func(t *testing.T) {
		inject := newContainer()

		const dbName = "main"
		const driverName = "some"

		sharedDep := dependency.NewShared(newDriver, dbName)

		if err := inject.Provide(sharedDep, Name(driverName)); err != nil {
			t.Error(err)
			return
		}

		usersRepo := dependency.New(newUserWithDriver, dependency.FromSharedWithName(driverName))
		if err := inject.Provide(usersRepo); err != nil {
			t.Error(err)
			return
		}

		todoRepo := dependency.New(newTodo, dependency.FromSharedWithName(driverName))
		if err := inject.Provide(todoRepo); err != nil {
			t.Error(err)
			return
		}

		type args struct {
			In
			UserRepo *user
			TodoRepo *todo
		}

		called := false

		invoker := func(in args) {
			assert.Equal(t, dbName, in.UserRepo.getDb().client())
			assert.Equal(t, dbName, in.TodoRepo.db.client())
			assert.Same(t, in.UserRepo.getDb(), in.TodoRepo.db)

			called = true
		}

		err := inject.Invoke(invoker)

		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("invoke with shared dependency by name that does not exist", func(t *testing.T) {
		inject := newContainer()

		const dbName = "main"
		const driverName = "some"

		sharedDep := dependency.NewShared(newDriver, dbName)

		if err := inject.Provide(sharedDep, Name(driverName)); err != nil {
			t.Error(err)
			return
		}

		usersRepo := dependency.New(newUserWithDriver, dependency.FromSharedWithName("algo"))
		if err := inject.Provide(usersRepo); err != nil {
			t.Error(err)
			return
		}

		type args struct {
			In
			UserRepo *user
		}

		invoker := func(in args) {}

		err := inject.Invoke(invoker)

		expErr := errors.New("inject: error resolving argument 0 for constructor func(inject.database) *inject.user: inject: can't solve shared dependency with name algo")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("invoke with shared dependency by name that return error", func(t *testing.T) {
		inject := newContainer()

		const dbName = "main"
		const driverName = "some"

		errDriver := errors.New("some")

		sharedDep := dependency.NewShared(func(_ string) (*driver, error) { return nil, errDriver }, dbName)

		if err := inject.Provide(sharedDep, Name(driverName)); err != nil {
			t.Error(err)
			return
		}

		usersRepo := dependency.New(newUserWithDriver, dependency.FromSharedWithName(driverName))
		if err := inject.Provide(usersRepo); err != nil {
			t.Error(err)
			return
		}

		type args struct {
			In
			UserRepo *user
		}

		invoker := func(in args) {}

		err := inject.Invoke(invoker)

		expErr := errors.New("inject: error resolving argument 0 for constructor func(inject.database) *inject.user: inject: can't solve shared dependency with name some: inject: error constructing `func(string) (*inject.driver, error)`: some")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})
}
