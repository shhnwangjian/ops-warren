package main

import "fmt"

func main() {
	errorPanic2()
}


/*
defer 表达式的函数如果在 panic 后面，则这个函数无法被执行
 */
func errorPanic() {
	panic("hello")
	defer func() {
		fmt.Println("word")
	}()
}

func errorPanic2() {
	defer func() {
		fmt.Println("hello")
	}()
	panic("word")
}
