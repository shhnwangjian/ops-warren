package main

import (
	"context"
	"fmt"
	"time"
)

// WithTimeout内部调用的就是WithDeadline
func main()  {
	ctx,cancel := context.WithTimeout(context.Background(),1 * time.Second)
	defer cancel()
	//go HelloHandle(ctx, 500*time.Millisecond) // 过期时间大于处理时间
	go HelloHandle(ctx, 2000*time.Millisecond) // 过期时间小于处理时间，上下文过期而终止执行HelloHandle
	select {
	case <- ctx.Done():
		fmt.Println("Hello WithTimeout demo ",ctx.Err())
	}

}

func HelloHandle(ctx context.Context, duration time.Duration)  {
	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	case <-time.After(duration):
		fmt.Println("process request with", duration)
	}
	fmt.Println("HelloHandle finish!")
}
