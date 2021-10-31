package main

import "fmt"

// filter 过滤掉不符合条件的元素
func filter(items []interface{}, fn func(index int, item interface{}) bool) []interface{} {
	var filteredItems []interface{}
	for index, value := range items {
		if fn(index, value) {
			filteredItems = append(filteredItems, value)
		}
	}
	return filteredItems
}

func main() {
	var nums []interface{}
	nums = append(nums, 1, 2, 3, 4, 5)
	evenNums := filter(nums, func(index int, num interface{}) bool { return num.(int)%2 == 0 })
	fmt.Printf("%d", evenNums[0].(int) + 2)
}