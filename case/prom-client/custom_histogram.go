package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*
直方图

带有 _bucket 后缀的计数器时间序列，使用 le（小于或等于） 标签指示该存储桶的上限，
具有上限的隐式存储桶 +Inf 也暴露于比最大配置的存储桶边界花费更长的时间的请求，
还包括使用后缀 _sum 累积总和和计数 _count 的指标，
这些时间序列中的每一个在概念上都是一个 counter 计数器（只能上升的单个值），只是它们是作为直方图的一部分创建的。
 */

func main() {
	// 创建带标签的 gauge 指标对象
	histogramMetrics := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "ops_custom_metrics_histogram",
			Help: "custom metrics(histogram)",
			// 所有 bucket 的下限都从 0 开始的，所以我们不需要明确配置每个 bucket 的下限，只需要配置上限即可
			Buckets: []float64{0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
	)

	histogramMetrics.Observe(0.911)
	// 注册到全局默认注册表中
	prometheus.MustRegister(histogramMetrics)

	go wait(histogramMetrics)

	// 暴露自定义的指标
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":9099", nil); err !=nil{
		log.Fatal(err)
	}
}

func wait(histogramMetrics prometheus.Histogram) {
	for {
		// 启动一个计时器
		timer := prometheus.NewTimer(histogramMetrics)

		// [...在应用中处理请求...]
		time.Sleep(300 * time.Millisecond)

		// 停止计时器并观察其持续时间，将其放进 requestDurations 的直方图指标中去
		timer.ObserveDuration()

		time.Sleep(1 * time.Second)
	}
}
