package inject

import (
	"fmt"
	"reflect"

	"github.com/Drafteame/inject/dependency"
)

// shared is a dependency injection container that holds global shared dependencies already built and ready to be
// injected.
type shared struct {
	solvedDeps       map[reflect.Type]any
	solvedDepsByName map[string]any

	deps       map[reflect.Type]dependency.Dependency
	depsByName map[string]dependency.Dependency
}

func newSharedContainer() *shared {
	return &shared{
		solvedDeps:       make(map[reflect.Type]any),
		solvedDepsByName: make(map[string]any),
		deps:             make(map[reflect.Type]dependency.Dependency),
		depsByName:       make(map[string]dependency.Dependency),
	}
}

func (sc *shared) GetByName(name string) (any, error) {
	if dep, ok := sc.solvedDepsByName[name]; ok {
		return dep, nil
	}

	builder, ok := sc.depsByName[name]
	if !ok {
		return nil, fmt.Errorf("inject: can't solve shared dependency with name %s", name)
	}

	res, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("inject: can't solve shared dependency with name %s: %v", name, err)
	}

	sc.solvedDepsByName[name] = res

	return res, nil
}

func (sc *shared) GetByType(t reflect.Type) (any, error) {
	if dep, ok := sc.solvedDeps[t]; ok {
		return dep, nil
	}

	builder, ok := sc.deps[t]
	if !ok {
		return nil, fmt.Errorf("inject: can't solve shared dependency by type %v", t)
	}

	res, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("inject: can't solve shared dependency by type %v: %v", t, err)
	}

	sc.solvedDeps[t] = res

	return res, nil
}
