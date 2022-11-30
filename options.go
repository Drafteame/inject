package inject

import (
	"fmt"
	"reflect"

	"github.com/Drafteame/inject/utils"
)

// injectOptions defines the possible options of an injection point.
type injectOptions struct {
	aliases map[reflect.Type]struct{}
	name    *string
}

// Option definition for an inject option
type Option func(*injectOptions) error

func newInjectOptions() injectOptions {
	return injectOptions{aliases: make(map[reflect.Type]struct{})}
}

// As It checks if the value passed to the `As` function is an interface or a pointer. If it's not, it returns an error.
// If it's either of them, then we add the type to a map of aliases for that particular injection point (the `opts`
// variable). We return nil if everything went well, otherwise we return an error message explaining what went wrong.
func As(alias any) Option {
	return func(opts *injectOptions) error {
		atype, err := utils.BuildAliasType(alias)
		if err != nil {
			return err
		}

		if opts.aliases == nil {
			opts.aliases = make(map[reflect.Type]struct{})
		}

		opts.aliases[atype] = struct{}{}

		return nil
	}
}

// Name It returns a function that takes an `injectOptions` pointer and returns an error. If the name is empty, it will
// return an error. Otherwise, it will set the name in the options struct and return nil (no error). The returned
// function is of type `Option`. The returned function can be used as a parameter to `inject.Get().Provide()`.
func Name(name string) Option {
	return func(opts *injectOptions) error {
		if name == "" {
			return fmt.Errorf("inject: name cannot be empty if provided")
		}

		opts.name = &name
		return nil
	}
}
