package sunfish

import "reflect"

func MakeNewElemFunction(sliceElementType reflect.Type) func() reflect.Value {
	var newElemFunc func() reflect.Value
	if sliceElementType.Kind() == reflect.Ptr {
		newElemFunc = func() reflect.Value {
			return reflect.New(sliceElementType.Elem())
		}
	} else {
		newElemFunc = func() reflect.Value {
			return reflect.New(sliceElementType)
		}
	}

	return newElemFunc
}

func MakeSliceValueSetFunc(sliceElementType reflect.Type, sliceValue reflect.Value) func(*reflect.Value) {
	var sliceValueSetFunc func(*reflect.Value)
	if sliceElementType.Kind() == reflect.Ptr {
		sliceValueSetFunc = func(newValue *reflect.Value) {
			sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(newValue.Interface())))
		}
	} else {
		sliceValueSetFunc = func(newValue *reflect.Value) {
			sliceValue.Set(reflect.Append(sliceValue, reflect.Indirect(reflect.ValueOf(newValue.Interface()))))
		}
	}
	return sliceValueSetFunc
}

