package container

import (
	"github.com/Drafteame/inject/dependency"
)

// Container is a dependency injection Container implementation
type Container struct {
	solvedDeps map[string]any
	deps       map[string]dependency.Dependency
}

// New creates a new instance of a Container.
func New() *Container {
	return &Container{
		solvedDeps: make(map[string]any),
		deps:       make(map[string]dependency.Dependency),
	}
}

// Flush WARNING: This function will delete all saved instances, solved and registered factories from the container.
// Do not use this method on production, and just use it on testing purposes.
func (c *Container) Flush() {
	c.solvedDeps = make(map[string]any)
	c.deps = make(map[string]dependency.Dependency)
}
