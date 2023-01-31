package container

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Drafteame/inject/dependency"
	"github.com/Drafteame/inject/types"
)

func TestContainer_Invoke(t *testing.T) {
	t.Run("invoke no dependency function", func(t *testing.T) {
		inject := New()
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

		inject := New()
		called := false
		depName := types.Symbol("test")

		type args struct {
			types.In
			UserService *user `inject:"name=test"`
		}

		invoker := func(in args) {
			if assert.NotNil(t, in.UserService) {
				assert.Equal(t, name, in.UserService.getName())
				assert.Equal(t, age, in.UserService.getAge())
			}

			called = true
		}

		userDep := dependency.New(newUser, name, age)

		if err := inject.Provide(depName, userDep); err != nil {
			t.Error(err)
			return
		}

		err := inject.Invoke(invoker)

		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("invoke with dependency injected to interface", func(t *testing.T) {
		const name = "John Smith"
		const age = 21
		const depName = "usersService"

		inject := New()
		called := false

		type args struct {
			types.In
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

		if err := inject.Provide(depName, userDep); err != nil {
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

		inject := New()
		called := false

		userDepName := types.Symbol("usersService")
		userDep := dependency.New(newUser, name, age)
		if err := inject.Provide(userDepName, userDep); err != nil {
			t.Error(err)
			return
		}

		type args struct {
			types.In
			UserService namer `inject:"name=test,optional"`
		}

		invoker := func(in args) {
			assert.Nil(t, in.UserService)

			called = true
		}

		err := inject.Invoke(invoker)

		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("invoke with empty inject tag dependency", func(t *testing.T) {
		const name = "John Smith"
		const age = 21

		inject := New()

		depName := types.Symbol("test")
		userDep := dependency.New(newUser, name, age)
		if err := inject.Provide(depName, userDep); err != nil {
			t.Error(err)
			return
		}

		type args struct {
			types.In
			UserService userer `inject:""`
		}

		invoker := func(in args) {}

		err := inject.Invoke(invoker)
		expErr := errors.New("inject: missing name tag of inject dependency on field `UserService`")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("invoke with nil invoker", func(t *testing.T) {
		inject := New()

		err := inject.Invoke(nil)

		expErr := fmt.Errorf("inject: can't invoke nil constructor")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("invoke with non-function invoker", func(t *testing.T) {
		inject := New()

		err := inject.Invoke(10)

		expErr := fmt.Errorf("inject: can't invoke a non-function constructor")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("invoke with non In embedded struct", func(t *testing.T) {
		inject := New()

		type args struct{}

		invoker := func(in args) {}

		err := inject.Invoke(invoker)

		expErr := fmt.Errorf("inject: struct doesn't embed `inject.In` struct")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("invoke with invoker returning non error value", func(t *testing.T) {
		inject := New()

		type args struct {
			types.In
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
		inject := New()

		type args struct {
			types.In
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
		inject := New()

		type args struct {
			types.In
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

	t.Run("invoke with error providing field with no tag", func(t *testing.T) {
		inject := New()

		type some interface {
			someMethod()
		}

		type args struct {
			types.In
			UserService some
		}

		invoker := func(in args) {}

		err := inject.Invoke(invoker)

		expErr := fmt.Errorf("inject: missing name tag of inject dependency on field `UserService`")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("invoke with error providing named dependency", func(t *testing.T) {
		inject := New()

		type args struct {
			types.In
			UserService *user `inject:"name=usersService"`
		}

		invoker := func(in args) {}

		err := inject.Invoke(invoker)

		expErr := fmt.Errorf("inject: no provided dependency of name `usersService`")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("invoke error resolving dependency three", func(t *testing.T) {
		inject := New()

		depName := types.Symbol("depName")
		dep := dependency.New(func() (*user, error) { return nil, errors.New("some") })
		if err := inject.Provide(depName, dep); err != nil {
			t.Error(err)
			return
		}

		type args struct {
			types.In
			UserService *user `inject:"name=depName"`
		}

		invoker := func(in args) {}

		err := inject.Invoke(invoker)

		expErr := fmt.Errorf("inject: error building dependency instance: inject: error constructing `func() (*container.user, error)`: some")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("invoke with shared dependency on multiple targets as singleton", func(t *testing.T) {
		inject := New()

		const driverName = "some"

		sharedDepName := types.Symbol("driver")
		sharedDep := dependency.NewSingleton(newDriver, driverName)
		if err := inject.Provide(sharedDepName, sharedDep); err != nil {
			t.Error(err)
			return
		}

		userRepoName := types.Symbol("userRepo")
		usersRepo := dependency.New(newUserWithDriver, dependency.Inject(sharedDepName))
		if err := inject.Provide(userRepoName, usersRepo); err != nil {
			t.Error(err)
			return
		}

		todoRepoName := types.Symbol("todoRepo")
		todoRepo := dependency.New(newTodo, dependency.Inject(sharedDepName))
		if err := inject.Provide(todoRepoName, todoRepo); err != nil {
			t.Error(err)
			return
		}

		type args struct {
			types.In
			UserRepo *user `inject:"name=userRepo"`
			TodoRepo *todo `inject:"name=todoRepo"`
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

	t.Run("invoke with shared dependency that return error", func(t *testing.T) {
		inject := New()

		const dbName = "main"

		errBuild := errors.New("some")

		sharedDepName := types.Symbol("driver")
		sharedDep := dependency.NewSingleton(func(dbname string) (*driver, error) { return nil, errBuild }, dbName)
		if err := inject.Provide(sharedDepName, sharedDep); err != nil {
			t.Error(err)
			return
		}

		usersRepoName := types.Symbol("usersRepo")
		usersRepo := dependency.New(newUserWithDriver, dependency.Inject(sharedDepName))
		if err := inject.Provide(usersRepoName, usersRepo); err != nil {
			t.Error(err)
			return
		}

		type args struct {
			types.In
			UserRepo *user `inject:"name=usersRepo"`
		}

		invoker := func(in args) {}

		err := inject.Invoke(invoker)
		expErr := errors.New("inject: error building dependency instance: inject: error resolving argument 0 for constructor func(container.database) *container.user: inject: error building dependency instance: inject: error constructing `func(string) (*container.driver, error)`: some")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})

	t.Run("invoke with shared dependency that does not exist", func(t *testing.T) {
		inject := New()

		const dbName = "main"

		sharedDepName := types.Symbol("driver")
		sharedDep := dependency.NewSingleton(newDriver, dbName)
		if err := inject.Provide(sharedDepName, sharedDep); err != nil {
			t.Error(err)
			return
		}

		usersRepoName := types.Symbol("usersRepo")
		usersRepo := dependency.New(newUserWithDriver, dependency.Inject("algo"))
		if err := inject.Provide(usersRepoName, usersRepo); err != nil {
			t.Error(err)
			return
		}

		type args struct {
			types.In
			UserRepo *user `inject:"name=usersRepo"`
		}

		invoker := func(in args) {}

		err := inject.Invoke(invoker)

		expErr := errors.New("inject: error building dependency instance: inject: error resolving argument 0 for constructor func(container.database) *container.user: inject: no provided dependency of name `algo`")

		assert.Error(t, err)
		assert.Equal(t, expErr, err)
	})
}
