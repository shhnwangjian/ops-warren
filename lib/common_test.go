package lib

import "testing"

func TestIf(t *testing.T) {
	a, b := 2, 3
	max := If(a > b, a, b).(int)
	println(max)
}
