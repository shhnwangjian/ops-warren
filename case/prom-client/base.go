package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*
go_*：以 go_ 为前缀的指标是关于 Go 运行时相关的指标，比如垃圾回收时间、goroutine 数量等，
这些都是 Go 客户端库特有的，其他语言的客户端库可能会暴露各自语言的其他运行时指标。

promhttp_*：来自 promhttp 工具包的相关指标，用于跟踪对指标请求的处理
 */

func main() {
	// Serve the default Prometheus metrics registry over HTTP on /metrics.
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":9099", nil); err !=nil{
		log.Fatal(err)
	}
}
