package main

import (
	"fmt"
)

/*
执行顺序是按调用 defer 语句的倒序执行
 */
func main() {
	defer func() {
		fmt.Println("first")
	}()

	defer func() {
		fmt.Println("second")
	}()

	fmt.Println("done")
}
