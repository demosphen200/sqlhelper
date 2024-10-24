package ex

func ValueOrPanic[T interface{}](value T, err error) T {
	if err != nil {
		panic(err.Error())
	} else {
		return value
	}
}

func ValueOrPanic2[T1 interface{}, T2 interface{}](value1 T1, value2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err.Error())
	} else {
		return value1, value2
	}
}
