package wrap

import (
	"fmt"
	"reflect"
)

func After(dst, src interface{}, afterFn func()) error {
	vDst := reflect.ValueOf(dst)
	vSrc := reflect.ValueOf(src)

	if vDst.Kind() != reflect.Ptr {
		return fmt.Errorf("First argument of WiretapWithArgs must be pointer of function")
	}

	vDst = vDst.Elem()
	if vDst.Kind() != reflect.Func {
		return fmt.Errorf("First argument of WiretapWithArgs must be pointer of function")
	}

	if vSrc.Kind() != reflect.Func {
		return fmt.Errorf("Second argument of WiretapWithArgs must be function")
	}

	if vDst.Type() != vSrc.Type() {
		return fmt.Errorf("Value of first argument must be same type function as second argument")
	}

	f := func(values []reflect.Value) []reflect.Value {
		defer afterFn()
		return vSrc.Call(values)
	}

	vDst.Set(reflect.MakeFunc(vDst.Type(), f))

	return nil
}
