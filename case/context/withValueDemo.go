package main

import (
	"context"
	"fmt"
)

type selfKey string

func main()  {
	ctx := context.WithValue(context.Background(), selfKey("hello"),"word")
	Get(ctx, selfKey("hello"))
	Get(ctx, selfKey("test"))
}

func Get(ctx context.Context, k selfKey)  {
	if v, ok := ctx.Value(k).(string); ok {
		fmt.Println(fmt.Sprintf("%v is value: %v", k, v))
		return
	}
	fmt.Println(fmt.Sprintf("%v no value", k))
}