package inject

import (
	"fmt"
	"reflect"

	"github.com/Drafteame/inject/container"
	"github.com/Drafteame/inject/dependency"
)

var injector Container

// Container represents a dependency container that should register factory methods and its dependency threes to be
// injected when
type Container interface {
	Provide(name string, dep dependency.Dependency) error
	Invoke(construct any) error
	Get(name string) (any, error)
	Flush()
}

// get return a global instance for the dependency injection container. If the container is nil, then it will initialize
// a new instance before returning the container.
func get() Container {
	if injector == nil {
		injector = container.New()
	}

	return injector
}

// New Return a new isolated instance for the dependency injection container. This instance is totally different from
// the global container and do not share any saved dependency three between each other.
func New() Container {
	return container.New()
}

// Provide Is a wrapper over the Provide function attached to the global container. It adds a new injection dependency
// to the container, getting the first result type of the constructor to associate the constructor on the injection
// dependency threes.
//
// This injection will be resolved and built on execution time when the `inject.Invoke(...)` or `inject.Get(name)`
// methods are called.
func Provide(name string, dep dependency.Dependency) error {
	return get().Provide(name, dep)
}

// Invoke Is the entry point to execute dependency injection resolution. It calls an invoker function that can
// receive or not a struct that embeds inject.In struct as input, and return an error or not (any other return field or
// type will be ignored on resolution). When invoker is called it will resolve the dependency threes of each field from
// the previously provided resources on Container.
func Invoke(construct any) error {
	return get().Invoke(construct)
}

// Get is a wrapper over the Get function attached to the global container. This function modify the return type of the
// resolved dependency, returned as `any` to the provided generic type `T`. If it can't be casted it will return an
// error.
func Get[T any](name string) (T, error) {
	instance, err := get().Get(name)
	if err != nil {
		aux := new(T)
		return *aux, err
	}

	cast, ok := instance.(T)
	if !ok {
		aux := new(T)
		axtype := reflect.TypeOf(*aux)
		return *aux, fmt.Errorf("inject: error casting instance of `%s` dependency to `%v`", name, axtype)
	}

	return cast, nil
}

// Flush WARNING: This function will delete all saved instances, solved and registered factories from the container.
// Do not use this method on production, and just use it on testing purposes.
func Flush() {
	get().Flush()
}
