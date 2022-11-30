package utils

import "reflect"

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
