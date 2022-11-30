package inject

import (
	"github.com/Drafteame/inject/dependency"
)

// Injector it's the representation of the dependency injection container.
type Injector interface {
	Provide(dependency.Dependency, ...Option) error
	Invoke(any) error
}

type Container interface {
	Injector

	GetByName(name string) (any, error)
	GetByType(dtype any) (any, error)
}

// Get return a global instance for the dependency injection container. If the container is nil, then it will initialize
// a new instance before returning the container.
func Get() Injector {
	if injector == nil {
		injector = newContainer()
	}

	return injector
}

// New Return a new isolated instance for the dependency injection container. This instance is totally different from
// the global container and do not share any saved dependency three between each other.
func New() Injector {
	return newContainer()
}
