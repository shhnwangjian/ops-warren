package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main()  {
	http.HandleFunc("/hello", Hello) // 设置访问的路由
	http.HandleFunc("/hello2", Hello2)

	if err := http.ListenAndServe(":9099", nil); err !=nil{
		log.Fatal(err)
	}
}

/*
Hello
假定请求需要耗时2s，在请求2s后返回，我们期望监控goroutine在打印2次Current request is in progress后即停止。
但运行发现，监控goroutine打印2次后，其仍不会结束，而会一直打印下去。
 */
func Hello(writer http.ResponseWriter, request *http.Request)  {
	fmt.Println(&request)

	go func() {
		for range time.Tick(time.Second) {
			fmt.Println("Current request is in progress(hello)")
		}
	}()

	time.Sleep(2 * time.Second)
	if _, err := writer.Write([]byte("Hi")) ; err !=nil{
		fmt.Println(err)
	}
}

/*
Hello2
context包可以提供一个请求从API请求边界到各goroutine的请求域数据传递、取消信号及截至时间等能力
 */
func Hello2(writer http.ResponseWriter, request *http.Request)  {
	fmt.Println(&request)

	go func() {
		for range time.Tick(time.Second) {
			select {
			case <- request.Context().Done():
				fmt.Println("request is outgoing(hello2)")
				return
			default:
				fmt.Println("Current request is in progress(hello2)")
			}
		}
	}()

	time.Sleep(2 * time.Second)
	if _, err := writer.Write([]byte("Hi 2")) ; err !=nil{
		fmt.Println(err)
	}
}