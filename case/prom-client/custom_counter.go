package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// 创建带标签的 gauge 指标对象
	counterMetrics := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ops_custom_metrics_counter",
			Help: "custom metrics(counter)",
		},
	)
	counterMetrics.Inc() // +1
	counterMetrics.Inc() // +1
	counterMetrics.Add(10) // +10

	// 注册到全局默认注册表中
	prometheus.MustRegister(counterMetrics)

	// 暴露自定义的指标
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":9099", nil); err !=nil{
		log.Fatal(err)
	}
}