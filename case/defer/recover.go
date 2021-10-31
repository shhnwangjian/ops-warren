package main

import "fmt"


/*
调用 A() 函数；
调用 B() 函数；
B() 中遇到 panic，立刻终止，不执行 panic 之后的代码；
执行 B() 中 defer 函数，遇到 recover 捕获错误，继续执行 defer 中代码，然后返回；
执行 A() 函数后续代码，最后执行 A() 中 defer 函数。
 */
func main() {
	A()
}

func A() {
	defer func() {
		fmt.Println("three")
	}()

	B()
	fmt.Println("继续执行A")
}

func B() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("捕获异常:", err)
		}
		fmt.Println("two")
	}()
	panic("one")
}