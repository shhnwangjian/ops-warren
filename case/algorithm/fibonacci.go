package main

import "fmt"

/*
描述
大家都知道斐波那契数列，现在要求输入一个正整数 n ，请你输出斐波那契数列的第 n 项。
斐波那契数列
数据范围：1\leq n\leq 391≤n≤39
要求：空间复杂度 O(1)O(1)，时间复杂度 O(n)O(n) ，本题也有时间复杂度 O(logn)O(logn) 的解法

输入描述：一个正整数n
返回值描述：输出一个正整数。
 */


func main() {
	res := Fibonacci(10)
	fmt.Println(res)
}

/*
Fibonacci
@param n int整型
@return int整型
 */
func Fibonacci(n int) int {
	// write code here
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}

	sum := 0
	n1 := 0
	n2 := 1

	for i:= 1; i < n; i++ {
		sum = n1 + n2
		n1 = n2
		n2 = sum
	}
	return sum
}