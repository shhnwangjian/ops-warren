package main

import (
	"context"
	"fmt"
	"time"
)

func main()  {
	ctx,cancel := context.WithCancel(context.Background())
	defer cancel()
	go cancelDemo(ctx)
	time.Sleep(6 * time.Second)
}

func cancelDemo(ctx context.Context)  {
	for range time.Tick(time.Second){
		select {
		case <- ctx.Done():
			return
		default:
			fmt.Println(fmt.Sprintf("hahahahahahahaha, %v", time.Now()))
		}
	}
}
