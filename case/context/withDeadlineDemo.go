package main

import (
	"context"
	"fmt"
	"time"
)

// WithDeadline 根据parent和deadline返回一个派生的Context。
// 如果parent存在过期时间，且已过期，则返回一个语义上等同于parent的派生Context。
// 当到达过期时间、或者调用CancelFunc函数关闭、或者关闭parent会使该函数返回的派生Context关闭。
func main()  {
	now := time.Now()
	later, _ := time.ParseDuration("1s")
	//later, _ := time.ParseDuration("3s")

	ctx, cancel := context.WithDeadline(context.Background(), now.Add(later))
	defer cancel()

	go deadlineDemo(ctx)

	time.Sleep(5 * time.Second)

}

func deadlineDemo(ctx context.Context)  {
	select {
	case <- ctx.Done():
		fmt.Println(ctx.Err())
	case <-time.After(2 * time.Second):
		fmt.Println("stop deadline demo")
	}
	fmt.Println("deadline demo finish!")
}
