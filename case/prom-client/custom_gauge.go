package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// 创建带标签的 gauge 指标对象
	gaugeMetrics := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ops_custom_metrics_gauge",
			Help: "custom metrics(gauge)",
		},
		// 指定标签名称
		[]string{"level", "degree"},
	)

	// 注册到全局默认注册表中
	prometheus.MustRegister(gaugeMetrics)

	// 针对不同标签值设置不同的指标值
	gaugeMetrics.WithLabelValues("3", "5").Set(27)
	gaugeMetrics.WithLabelValues("low", "10").Set(25.3)

	// 暴露自定义的指标
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":9099", nil); err !=nil{
		log.Fatal(err)
	}
}