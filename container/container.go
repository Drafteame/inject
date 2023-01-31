package container

import (
	"github.com/Drafteame/inject/dependency"
	"github.com/Drafteame/inject/types"
)

// Container is a dependency injection Container implementation
type Container struct {
	solvedDeps map[types.Symbol]any
	deps       map[types.Symbol]dependency.Dependency
}

// New creates a new instance of a Container.
func New() *Container {
	return &Container{
		solvedDeps: make(map[types.Symbol]any),
		deps:       make(map[types.Symbol]dependency.Dependency),
	}
}

// Flush WARNING: This function will delete all saved instances, solved and registered factories from the container.
// Do not use this method on production, and just use it on testing purposes.
func (c *Container) Flush() {
	c.solvedDeps = make(map[types.Symbol]any)
	c.deps = make(map[types.Symbol]dependency.Dependency)
}
