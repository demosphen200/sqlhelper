package utils

func Xor(value1 bool, value2 bool) bool {
	return (value1 && !value2) || (!value1 && value2)
}
