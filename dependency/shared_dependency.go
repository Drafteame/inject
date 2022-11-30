package dependency

import (
	"fmt"
	"reflect"

	"github.com/Drafteame/inject/utils"
)

const (
	sourceDeps      = "deps"
	sourceNamedDeps = "named-deps"
)

type Shared struct {
	depName         string
	depType         reflect.Type
	source          string
	sharedContainer SharedContainer
}

var _ Builder = &Shared{}

func FromSharedWithName(name string) Shared {
	return Shared{
		depName: name,
		source:  sourceNamedDeps,
	}
}

func FromSharedWithType(alias any) Shared {
	atype, err := utils.BuildAliasType(alias)
	if err != nil {
		panic(err)
	}

	return Shared{
		depType: atype,
		source:  sourceDeps,
	}
}

func (s Shared) Build() (any, error) {
	if s.sharedContainer == nil {
		return nil, fmt.Errorf("inject: [internal-error] no shared container provided")
	}

	if s.source == sourceDeps {
		return s.sharedContainer.GetByType(s.depType)
	}

	return s.sharedContainer.GetByName(s.depName)
}

func (s Shared) IsShared() bool {
	return false
}

func (s Shared) WithSharedContainer(c SharedContainer) Builder {
	s.sharedContainer = c
	return s
}
