package container

import (
	"fmt"

	"github.com/Drafteame/inject/dependency"
	"github.com/Drafteame/inject/types"
	"github.com/Drafteame/inject/utils"
)

// Provide It adds a new injection dependency to the Container, getting the first result type of the constructor to
// associate the constructor on the injection dependency threes, e.g:
//
// inject.get().Provide(dependency.New(callback, arg1, arg2), inject.As(new(someInterface)))
//
// This injection will be resolved and built on execution time when the `inject.get().Invoke(...)` method is called.
func (c *Container) Provide(name types.Symbol, dep dependency.Dependency) error {
	var err error

	if rt := utils.GetFirstReturnType(dep.Factory); rt == nil {
		return fmt.Errorf("inject: dependency factory should return at least one return type: %s", dep.String())
	}

	c.deps, err = c.provide(c.deps, name, dep)
	if err != nil {
		return err
	}

	return nil
}

// provide If the name option is not set, it returns the Container and nil. If the Container is nil, it creates a
// new one. It checks if there's already a dependency with that name in the Container and returns an error if so. It
// adds the dependency to the Container using its name as key and returns it along with nil (no error).
func (c *Container) provide(container map[types.Symbol]dependency.Dependency, name types.Symbol, dep dependency.Dependency) (map[types.Symbol]dependency.Dependency, error) {
	if name == "" {
		return container, fmt.Errorf("inject: dependency name cannot be empty")
	}

	if container == nil {
		container = make(map[types.Symbol]dependency.Dependency)
	}

	if _, ok := container[name]; ok {
		return container, fmt.Errorf("inject: duplicated dependency name `%s`", name)
	}

	container[name] = dep

	return container, nil
}
