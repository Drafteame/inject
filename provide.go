package inject

import (
	"fmt"
	"reflect"

	"github.com/Drafteame/inject/dependency"
	"github.com/Drafteame/inject/utils"
)

// Provide It adds a new injection dependency to the container, getting the first result type of the constructor to
// associate the constructor on the injection dependency threes, e.g:
//
// inject.Get().Provide(dependency.New(callback, arg1, arg2), inject.As(new(someInterface)))
//
// This injection will be resolved and built on execution time when the `inject.Get().Invoke(...)` method is called.
func (c *container) Provide(dep dependency.Dependency, opts ...Option) error {
	if dep.IsShared() {
		return c.provideShared(dep, opts...)
	}

	return c.provideIsolated(dep, opts...)
}

func (c *container) provideShared(dep dependency.Dependency, opts ...Option) error {
	options := newInjectOptions()

	for _, option := range opts {
		if err := option(&options); err != nil {
			return err
		}
	}

	rtype := utils.GetFirstReturnType(dep.Constructor)
	if rtype == nil {
		return fmt.Errorf("inject: can't provide a dependency constructor with no return types: %s", dep.String())
	}

	options.aliases[rtype] = struct{}{}

	var err error

	c.shared.deps = c.provideAliases(c.shared.deps, options, dep)

	c.shared.depsByName, err = c.provideName(c.shared.depsByName, options, dep)
	if err != nil {
		return err
	}

	return nil
}

// provideIsolated It creates a new `injectOptions` struct. It iterates over the options and applies them to the
// `injectOptions` struct. It gets the first return type of the dependency constructor function (the one that will be
// used for injection). If there is no return type, it returns an error saying that we can't provide a dependency with
// no return types (this is not possible in Go). It adds this first return type to the aliases list of our
// `injectOptions`. This means that when we inject this dependency, we will also inject all its aliases (if they are
// registered in our container). Then it calls two functions: `provideAliases` and `provideName`. These functions are
// responsible for adding dependencies to their respective lists: by alias or by name (we'll see how these work later
// on).
func (c *container) provideIsolated(dep dependency.Dependency, opts ...Option) error {
	options := newInjectOptions()

	for _, option := range opts {
		if err := option(&options); err != nil {
			return err
		}
	}

	rtype := utils.GetFirstReturnType(dep.Constructor)
	if rtype == nil {
		return fmt.Errorf("inject: can't provide a dependency constructor with no return types: %s", dep.String())
	}

	options.aliases[rtype] = struct{}{}

	var err error

	c.deps = c.provideAliases(c.deps, options, dep)

	c.depsByName, err = c.provideName(c.depsByName, options, dep)
	if err != nil {
		return err
	}

	return nil
}

// provideAliases It checks if the `opts.aliases` is empty, and if it is, it returns the container as-is. If the
// container is nil, it creates a new one (this should never happen). For each alias in `opts`, we add an entry to the
// container with that alias as key and dep as value. We return the updated container and no error (nil).
func (c *container) provideAliases(container map[reflect.Type]dependency.Dependency, opts injectOptions, dep dependency.Dependency) map[reflect.Type]dependency.Dependency {
	if container == nil {
		container = make(map[reflect.Type]dependency.Dependency)
	}

	for alias, _val := range opts.aliases {
		_ = _val
		container[alias] = dep
	}

	return container
}

// provideName If the name option is not set, it returns the container and nil. If the container is nil, it creates a
// new one. It checks if there's already a dependency with that name in the container and returns an error if so. It
// adds the dependency to the container using its name as key and returns it along with nil (no error).
func (c *container) provideName(container map[string]dependency.Dependency, opts injectOptions, dep dependency.Dependency) (map[string]dependency.Dependency, error) {
	if opts.name == nil {
		return container, nil
	}

	if container == nil {
		container = make(map[string]dependency.Dependency)
	}

	if _, ok := container[*opts.name]; ok {
		return container, fmt.Errorf("inject: duplicated dependency name `%s`", *opts.name)
	}

	container[*opts.name] = dep

	return container, nil
}
