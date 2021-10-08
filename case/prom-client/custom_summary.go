package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*
摘要与直方图非常类似，只是我们需要指定要跟踪的 quantiles 分位数值，而不需要处理 bucket 桶，
比如我们想要跟踪 HTTP 请求延迟的第 50、90 和 99 个百分位
*/

func main() {
	// 创建带标签的 gauge 指标对象
	summaryMetrics := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "ops_custom_metrics_summary",
			Help: "custom metrics(summary)",
			Objectives: map[float64]float64{
				0.5: 0.05,   // 第50个百分位数，最大绝对误差为0.05。
				0.9: 0.01,   // 第90个百分位数，最大绝对误差为0.01。
				0.99: 0.001, // 第90个百分位数，最大绝对误差为0.001。
			},
		},
	)

	summaryMetrics.Observe(0.911)
	summaryMetrics.Observe(0.3)
	// 注册到全局默认注册表中
	prometheus.MustRegister(summaryMetrics)


	// 暴露自定义的指标
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":9099", nil); err !=nil{
		log.Fatal(err)
	}
}
