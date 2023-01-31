package dependency

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Drafteame/inject/types"
)

//go:generate mockery --name=Builder --filename=builder.go --structname=Builder --output=mocks --outpkg=mocks
//go:generate mockery --name=Container --filename=container.go --structname=Container --output=mocks --outpkg=mocks

// Builder definition for a dependency that should be build on injection time.
type Builder interface {
	Build() (any, error)
}

// Container is a container that holds global dependencies.
type Container interface {
	Get(name types.Symbol) (any, error)
}

// Dependency implementation of dependency.
type Dependency struct {
	Factory   any
	Args      []any
	Singleton bool
	container Container
}

// New Create a new Dependency struct to build injection. Factory is a function with one of the next
// supported signs:
//   - func()
//   - func() any
//   - func() (any, error)
//   - func(any-arguments)
//   - func(any-arguments) any
//   - func(any-arguments) (any, error)
func New(constructor any, args ...any) Dependency {
	return Dependency{
		Factory: constructor,
		Args:    args,
	}
}

// NewSingleton Create a new Dependency struct to build injection but marking that will be a shared dependency to provide.
func NewSingleton(constructor any, args ...any) Dependency {
	return Dependency{
		Factory:   constructor,
		Args:      args,
		Singleton: true,
	}
}

// IsSingleton returns true if the current dependency will be treated as a shared dependency.
func (d Dependency) IsSingleton() bool { return d.Singleton }

// SetContainer add shared container to the dependency object in order to resolve shared arguments in the
// dependency three.
func (d Dependency) SetContainer(sc Container) Dependency {
	d.container = sc
	return d
}

// Build It validates the constructor and gets its type. It gets the arguments values for the constructor. It calls the
// constructor with those arguments using reflection (`reflect` package). Finally, it returns a value and an error if
// any of them is not nil (the error can be returned by one of the dependencies).
func (d Dependency) Build() (any, error) {
	ctype, err := d.validateAndGetReflectType()
	if err != nil {
		return nil, err
	}

	args, err := d.getArgsValues(ctype)
	if err != nil {
		return nil, err
	}

	res := reflect.ValueOf(d.Factory).Call(args)

	arg, err := d.getValueAndError(res)
	if err != nil {
		return nil, fmt.Errorf("inject: error constructing `%v`: %v", ctype, err)
	}

	return arg, nil
}

func (d Dependency) String() string {
	ctype := reflect.TypeOf(d.Factory)
	return fmt.Sprintf("dependency.Dependency{Factory: %v, Args: %v}", ctype, d.Args)
}

// validateAndGetReflectType It checks if the constructor is a function. It checks if the number of arguments provided
// to the constructor matches the number of arguments expected by it. If everything is ok, it returns a `reflect.Type`
// object that represents the type of our constructor function.
func (d Dependency) validateAndGetReflectType() (reflect.Type, error) {
	ctype := reflect.TypeOf(d.Factory)
	if ctype == nil {
		return nil, errors.New("inject: can't build an untyped nil")
	}

	if ctype.Kind() != reflect.Func {
		return nil, fmt.Errorf("inject: must provide constructor function, got `%v`", ctype)
	}

	argsLen := ctype.NumIn()
	if argsLen != len(d.Args) {
		return nil, fmt.Errorf("inject: invalid argument length for constructor `%v`, got %v (need %v)", ctype, len(d.Args), argsLen)
	}

	return ctype, nil
}

// getArgsValues It creates a slice of reflect.Value with the length of the number of arguments that we want to pass to
// the constructor. Then it iterates over all arguments and checks if they are assignable to the type that is expected
// by the constructor (the type is taken from ctype). If they are not assignable, then an error is returned, otherwise
// it adds them to values slice as reflect.Value objects and returns them at last along with nil error value
// (if everything went well).
func (d Dependency) getArgsValues(ctype reflect.Type) ([]reflect.Value, error) {
	args, err := d.resolveArguments(ctype)
	if err != nil {
		return nil, err
	}

	values := make([]reflect.Value, len(args))

	for i := 0; i < len(args); i++ {
		targ := ctype.In(i)

		if (targ.Kind() == reflect.Interface || targ.Kind() == reflect.Ptr) && args[i] == nil {
			values[i] = reflect.Zero(targ)
			continue
		}

		xt := reflect.TypeOf(args[i])

		if !xt.AssignableTo(targ) {
			return nil, fmt.Errorf("inject: using %s as type %s on constructor `%v`", xt.String(), targ.String(), ctype)
		}

		values[i] = reflect.ValueOf(args[i])
	}

	return values, nil
}

// resolveArguments It creates a slice of `any` type with the length of the number of arguments. For each argument, it
// normalizes it and builds it using the builder. If there is an error, we return an error message that contains
// information about which argument failed and what constructor was used (we will see how this works in a moment).
// Otherwise, we add the result to our slice and continue with the next argument until all arguments are resolved or an
// error occurs. Finally, we return our slice of arguments or an error if one occurred during resolution.
func (d Dependency) resolveArguments(ctype reflect.Type) ([]any, error) {
	args := make([]any, len(d.Args))

	var res any
	var err error

	for i := 0; i < len(d.Args); i++ {
		switch d.Args[i].(type) {
		case Injectable:
			arg := d.Args[i].(Injectable).SetContainer(d.container)
			res, err = d.resolveArgument(i, arg, ctype)
		case Dependency:
			arg := d.Args[i].(Dependency).SetContainer(d.container)
			res, err = d.resolveArgument(i, arg, ctype)
		default:
			arg := d.normalizeArgument(d.Args[i]).SetContainer(d.container)
			res, err = d.resolveArgument(i, arg, ctype)
		}

		if err != nil {
			return nil, err
		}

		args[i] = res
	}

	return args, nil
}

func (d Dependency) resolveArgument(index int, builder Builder, ctype reflect.Type) (any, error) {
	errMsg := "inject: error resolving argument %d for constructor %v: %v"
	res, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf(errMsg, index, ctype, err)
	}

	return res, nil
}

// getValueAndError If the constructor returns no values, we return `nil` and `nil`. If the constructor returns one
// value, we return that value and `nil`. If the constructor returns more than one value, we take the first as a result
// and last as an error. We check if last argument is an error (if it's not nil). We return result and error (or just
// nil if there was no error).
func (d Dependency) getValueAndError(res []reflect.Value) (any, error) {
	if len(res) == 0 {
		return nil, nil
	}

	if len(res) == 1 {
		return res[0].Interface(), nil
	}

	argValue := res[0].Interface()
	argErr := res[len(res)-1].Interface()

	var err error

	if argErr != nil {
		var ok bool

		err, ok = argErr.(error)
		if !ok {
			return nil, fmt.Errorf("inject: last result argument of the constructor is not an error")
		}
	}

	return argValue, err
}

// normalizeArgument If the argument is nil, we return a new dependency that returns nil. If the argument is already a
// builder, we return it as-is. If the argument is a function, we wrap it in a new dependency and return that instead.
// Otherwise, we wrap the value in a new dependency and return that instead (this will be used for constants).
func (d Dependency) normalizeArgument(arg any) Dependency {
	if arg == nil {
		return New(func() any { return arg })
	}

	if atype := reflect.TypeOf(arg); atype.Kind() == reflect.Func {
		return New(arg)
	}

	return New(func() any { return arg })
}
