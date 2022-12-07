package container

import (
	"fmt"
	"reflect"

	"github.com/Drafteame/inject/types"
	"github.com/Drafteame/inject/utils"
)

// Invoke Is the entry point to execute dependency injection resolution. It calls an invoker function that can
// receive or not a struct that embeds inject.In struct as input, and return an error or not (any other return field or
// type will be ignored on resolution). When invoker is called it will resolve the dependency threes of each field from
// the previously provided resources on Container.
func (c *Container) Invoke(construct any) error {
	if construct == nil {
		return fmt.Errorf("inject: can't invoke nil constructor")
	}

	ctype := reflect.TypeOf(construct)

	if ctype.Kind() != reflect.Func {
		return fmt.Errorf("inject: can't invoke a non-function constructor")
	}

	args, err := c.getInDeps(ctype)
	if err != nil {
		return err
	}

	res := reflect.ValueOf(construct).Call(args)

	return getResponseError(ctype, res)
}

// getResponseError It gets the type of the function that is being called. It checks if the function has an error as
// output parameter. If it does, it creates a new instance of `definitions.Error` and sets its value to the error
// returned by the function call (the index is calculated in step 2). Finally, it returns this error as an interface
// that can be casted from `definitions.Error` to `error`.
func getResponseError(ctype reflect.Type, res []reflect.Value) error {
	index, hasErr := utils.WhereErrorOut(ctype)
	if !hasErr {
		return nil
	}

	errInt := reflect.TypeOf(new(types.Error)).Elem()
	err := reflect.New(errInt)
	err.Elem().Set(res[index])

	return *err.Interface().(*types.Error)
}

// getInDeps It creates a slice of reflect.Value with the size of the number of input parameters. For each input
// parameter, it creates a new `reflect.Value` using `reflect.New`. Then it calls `buildInStruct` to build the struct
// and set its fields. If the type or the input struct is not a pointer, we need to get its value using `Elem()` method.
// We add this value to our slice of values and return it at the end.
func (c *Container) getInDeps(ctype reflect.Type) ([]reflect.Value, error) {
	values := make([]reflect.Value, ctype.NumIn())

	for i := 0; i < ctype.NumIn(); i++ {
		newArg := reflect.New(ctype.In(i))

		if err := types.BuildIn(c, newArg); err != nil {
			return nil, err
		}

		if ctype.In(i).Kind() != reflect.Ptr {
			newArg = newArg.Elem()
		}

		values[i] = newArg
	}

	return values, nil
}
