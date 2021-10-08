package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*
自定义一个的 gauge 指标
1、使用 prometheus.NewGauge() 函数创建了一个自定义的 gauge 指标对象，指标名称为 ops_custom_metrics_gauge，并添加了一个注释信息。
2、使用 prometheus.MustRegister() 函数在默认的注册表中注册了这个 gauge 指标对象。
3、通过调用 Set() 方法将 gauge 指标的值设置为 911。
4、通过 HTTP 暴露默认的注册表。

需要注意的是除了 prometheus.MustRegister() 函数之外还有一个 prometheus.Register() 函数，
一般在 golang 中我们会将 Mustxxx 开头的函数定义为必须满足条件的函数，
如果不满足会返回一个 panic 而不是一个 error 操作，所以如果这里不能正常注册的话会抛出一个 panic。
 */

func main() {
	// 创建一个没有任何 label 标签的 gauge 指标
	gaugeMetrics := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ops_custom_metrics_gauge",
		Help: "custom metrics(gauge)",
	})

	// 在默认的注册表中注册该指标
	prometheus.MustRegister(gaugeMetrics)

	//if err := prometheus.Register(gaugeMetrics); err !=nil {
	//	log.Fatal(err)
	//}

	// 设置 gauge 的值为 911
	gaugeMetrics.Set(911)

	// 暴露指标
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":9099", nil); err !=nil{
		log.Fatal(err)
	}
}