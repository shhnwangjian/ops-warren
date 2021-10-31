package main

import "fmt"

/*
描述
写出一个程序，接受一个字符串，然后输出该字符串反转后的字符串。（字符串长度不超过1000）

数据范围： 0 \le n \le 10000≤n≤1000
要求：空间复杂度 O(n)O(n)，时间复杂度 O(n)O(n)

输入："abcd"
返回值："dcba"
 */

func main() {
	res := solveByte("abcd")
	fmt.Println(res)
	res = solveStr("abcd")
	fmt.Println(res)
}

func solveByte(str string) string {
	b := []byte(str)
	count := len(b)
	for i := 0; i < count/2; i++ {
		t := b[i]
		b[i] = str[count-i-1]
		b[count-i-1] = t
	}
	return string(b)
}
func solveStr(str string) string {
	var res string
	count := len(str)
	for i := 0; i < count; i++ {
		tmp := str[i]
		res =  string(tmp) + res
	}
	return res
}