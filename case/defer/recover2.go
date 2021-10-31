package main

import "fmt"

/*
调用 C() 函数；
调用 D() 函数；
D() 中遇到 panic，立刻终止，不执行 panic 之后的代码；
执行 D() 中 defer 函数，由于没有 recover，则将 panic 抛到 C() 中；
C() 收到 panic 则不会执行后续代码，直接执行 defer 函数；
defer 中捕获 D() 抛出的异常，然后继续执行，最后退出。
 */
func main() {
	C()
}

func C() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("捕获异常:", err)
		}
		fmt.Println("three")
	}()

	D()
	fmt.Println("继续执行C")
}

func D() {
	defer func() {
		fmt.Println("one")
	}()
	panic("two")
}