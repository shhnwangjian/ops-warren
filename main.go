package main

import (
	"fmt"
	"math/rand"
	"ops-warren/balance"
	"time"
)

func main() {
	balanceTest() // 负载算法测试（全部）
}

func balanceTest() {
	fmt.Println("START-----------------------------------")
	var instanceList []*balance.Instance
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 5; i++ {
		host := fmt.Sprintf("10.10.%d.%d", rand.Intn(255), i)
		w := rand.Intn(10)
		fmt.Println(host, w)
		one := balance.NewInstance(host, 8080, int64(w))
		instanceList = append(instanceList, one)
	}
	fmt.Println("-----------------------------------")
	var balanceNames = []string{"hash", "random", "roundrobin", "weight_roundrobin", "shuffle"}

	for _, name := range balanceNames {
		startTime := time.Now().UnixNano()
		for i := 0; i < 10000; i++ {
			_, err := balance.DoBalance(name, instanceList)
			if err != nil {
				fmt.Println("Do balance err:", err)
				continue
			}
		}
		endTime := time.Now().UnixNano()
		fmt.Println("name: ", name, "cost time: ", (endTime-startTime)/1000) // 微秒
		for _, inst := range instanceList {
			if name == "weight_roundrobin" {
				fmt.Println(inst.GetResult(), ";weight: ", inst.Weight)
			} else {
				fmt.Println(inst.GetResult())
			}
			inst.CallNums = 0
		}
		fmt.Println("-----------------------------------")
	}
}
