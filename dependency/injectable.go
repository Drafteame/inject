package dependency

import (
	"fmt"
)

// Injectable is a type of dependency that is not a dependency three itself, but also is a reference to other dependency
// three, stored on the container. This Dependency will be accessed by its associated name on the container.
type Injectable struct {
	name      string
	container Container
}

var _ Builder = &Injectable{}

func Inject(name string) Injectable {
	return Injectable{
		name: name,
	}
}

func (s Injectable) Build() (any, error) {
	if s.container == nil {
		return nil, fmt.Errorf("inject: [internal-error] no container provided")
	}

	return s.container.Get(s.name)
}

func (s Injectable) IsSingleton() bool {
	return false
}

func (s Injectable) SetContainer(c Container) Builder {
	s.container = c
	return s
}
