package utils

func EuclideanModule(value int, module int) int {
	return ((value % module) + module) % module
}
