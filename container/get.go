package container

import (
	"fmt"

	"github.com/Drafteame/inject/dependency"
	"github.com/Drafteame/inject/types"
)

// Get returns a dependency instance and a possible build error by the associated name on that dependency. The instance
// type will depend on the dependency configuration, if it was marked as a singleton or not. If it was, the builder will
// try to return a previously created instance of that dependency instead of just create a new instance.
func (c *Container) Get(name types.Symbol) (any, error) {
	dep, ok := c.deps[name]
	if !ok {
		return nil, fmt.Errorf("inject: no provided dependency of name `%s`", name)
	}

	if dep.IsSingleton() {
		return c.getSingleton(name, dep)
	}

	return c.getInstance(dep)
}

func (c *Container) getSingleton(name types.Symbol, dep dependency.Dependency) (any, error) {
	val, ok := c.solvedDeps[name]
	if ok {
		return val, nil
	}

	val, err := c.getInstance(dep)
	if err != nil {
		return nil, err
	}

	c.solvedDeps[name] = val

	return val, nil
}

func (c *Container) getInstance(dep dependency.Dependency) (any, error) {
	val, err := dep.SetContainer(c).Build()
	if err != nil {
		return nil, fmt.Errorf("inject: error building dependency instance: %v", err)
	}

	return val, nil
}
