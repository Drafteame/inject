package utils

import "reflect"

// WhereErrorOut it tries to find from the reflected type of function callback, if this call back has any return error
// value and if it does, what's the index of the return value where the error comes from. If the function has more than
// one error value as return types, it will return the index of the first error founded.
func WhereErrorOut(ctype reflect.Type) (int, bool) {
	errInt := reflect.TypeOf((*error)(nil)).Elem()

	if ctype.Kind() != reflect.Func {
		return 0, false
	}

	nout := ctype.NumOut()

	if nout == 0 {
		return 0, false
	}

	if nout == 1 && ctype.Out(0).Kind() == reflect.Interface && ctype.Out(0).Implements(errInt) {
		return 0, true
	}

	for i := 0; i < nout; i++ {
		out := ctype.Out(i)

		if out.Kind() == reflect.Interface && out.Implements(errInt) {
			return i, true
		}
	}

	return 0, false
}
