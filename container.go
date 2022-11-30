package inject

import (
	"reflect"

	"github.com/Drafteame/inject/dependency"
)

var injector *container

// container is a dependency injection container implementation
type container struct {
	solvedDeps       map[reflect.Type]any
	solvedDepsByName map[string]any

	deps       map[reflect.Type]dependency.Dependency
	depsByName map[string]dependency.Dependency
	shared     *shared
}

func newContainer() *container {
	return &container{
		solvedDeps:       make(map[reflect.Type]any),
		solvedDepsByName: make(map[string]any),
		deps:             make(map[reflect.Type]dependency.Dependency),
		depsByName:       make(map[string]dependency.Dependency),
		shared:           newSharedContainer(),
	}
}
