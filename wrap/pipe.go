package wrap

import (
	"fmt"
	"reflect"
)

func Pipe(dst interface{}, fnList ...interface{}) error {
	if len(fnList) == 0 {
		return nil
	}

	vDst := reflect.ValueOf(dst)
	if vDst.Kind() != reflect.Ptr {
		return fmt.Errorf("First argument of WiretapWithArgs must be pointer of function")
	}

	vDst = vDst.Elem()
	if vDst.Kind() != reflect.Func {
		return fmt.Errorf("First argument of WiretapWithArgs must be pointer of function")
	}

	vFn := reflect.ValueOf(fnList[0])
	if vFn.Kind() != reflect.Func {
		return fmt.Errorf("arguments after the second must be function")
	}

	tDst := vDst.Type()
	tFn := vFn.Type()
	dstInNum, fnInNum := tDst.NumIn(), tFn.NumIn()
	if dstInNum != fnInNum {
		return fmt.Errorf("second argument must have arguments type same as arguments type of value of first argument")
	}

	for i := 0; i < dstInNum; i++ {
		if tDst.In(i) != tFn.In(i) {
			return fmt.Errorf("second argument must have arguments type same as arguments type of value of first argument")
		}
	}

	if len(fnList) == 1 {
		dstOutNum, fnOutNum := tDst.NumOut(), tFn.NumOut()
		if dstOutNum != fnOutNum {
			return fmt.Errorf("last argument must have returned values type same as returned values type of value of first argument")
		}

		for i := 0; i < dstOutNum; i++ {
			if tDst.Out(i) != tFn.Out(i) {
				return fmt.Errorf("last argument must have returned values type same as returned values type of value of first argument")
			}
		}

		vDst.Set(reflect.MakeFunc(vDst.Type(), func(values []reflect.Value) []reflect.Value {
			return reflect.ValueOf(fnList[0]).Call(values)
		}))
		return nil
	}

	for _, fn := range fnList[1:] {
		_vFn := reflect.ValueOf(fn)
		if _vFn.Kind() != reflect.Func {
			return fmt.Errorf("arguments after the second must be function")
		}

		_tFn := _vFn.Type()
		numOut, numIn := tFn.NumOut(), _tFn.NumIn()
		if numOut != numIn {
			return fmt.Errorf("all the arguments after the second must have returned values type same as arguments type of next function")
		}

		for i := 0; i < numOut; i++ {
			if tFn.Out(i) != _tFn.In(i) {
				return fmt.Errorf("all the arguments after the second must have returned values type same as arguments type of next function")
			}
		}
		tFn = _tFn
	}

	dstOutNum, fnOutNum := tDst.NumOut(), tFn.NumOut()
	if dstOutNum != fnOutNum {
		return fmt.Errorf("last argument must have returned values type same as returned values type of value of first argument")
	}

	for i := 0; i < dstOutNum; i++ {
		if tDst.Out(i) != tFn.Out(i) {
			return fmt.Errorf("last argument must have returned values type same as returned values type of value of first argument")
		}
	}

	f := func(values []reflect.Value) []reflect.Value {
		args := values
		for _, fn := range fnList {
			args = reflect.ValueOf(fn).Call(args)
		}
		return args
	}

	vDst.Set(reflect.MakeFunc(vDst.Type(), f))

	return nil
}
