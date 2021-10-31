package main

import (
	"fmt"
	"sort"
)

/*
描述
给定一个长度为 n 的可能有重复值的数组，找出其中不去重的最小的 k 个数。例如数组元素是4,5,1,6,2,7,3,8这8个数字，则最小的4个数字是1,2,3,4(任意顺序皆可)。
数据范围：0\le k,n \le 100000≤k,n≤10000，数组中每个数的大小0 \le val \le 10000≤val≤1000
要求：空间复杂度 O(n)O(n) ，时间复杂度 O(nlogn)O(nlogn)
 */

/*
GetLeastNumbersSolution
 * @param input int整型一维数组
 * @param k int整型
 * @return int整型一维数组
 */
func GetLeastNumbersSolution(input []int ,  k int ) []int {
	// write code here
	if len(input) == 0 {
		return []int{}
	}
	sort.Ints(input)
	return input[:k]
}

func main() {
	fmt.Println(GetLeastNumbersSolution([]int{4,5,1,6,2,7,3,8},4))
	fmt.Println(GetLeastNumbersSolution2([]int{4,5,1,6,2,7,3,8},4))
}

// GetLeastNumbersSolution2 修改快排，平均nlogn,最坏On^2 ; 空间平均logn，最坏On
func GetLeastNumbersSolution2(input []int, k int) []int {
	if k == 0 {
		return []int{}
	}
	quickSort(input, 0 , len(input) - 1)
	return input[:k]
}

func quickSort(nums []int, left, right int) {
	if left > right {
		return
	}

	i, j, base := left, right, nums[left]               //优化点，随机选基准

	for i < j {
		for nums[j] >= base && i < j {
			j--
		}
		for nums[i] <= base && i < j {
			i++
		}
		nums[i], nums[j] = nums[j], nums[i]
	}
	nums[i], nums[left] = nums[left], nums[i]

	quickSort(nums, left, i - 1)
	quickSort(nums, i + 1, right)
}